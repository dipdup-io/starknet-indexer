package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tester",
		Short: "Tester for indexer database",
	}
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})
	configPath := rootCmd.PersistentFlags().StringP("config", "c", "dipdup.yml", "path to YAML config file")
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("command line execute")
		return
	}
	if err := rootCmd.MarkFlagRequired("config"); err != nil {
		log.Panic().Err(err).Msg("config command line arg is required")
		return
	}

	var cfg Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelInfoValue
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Panic().Err(err).Msg("parsing log level")
		return
	}
	zerolog.SetGlobalLevel(logLevel)

	ctx, cancel := context.WithCancel(context.Background())

	postgres, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		log.Panic().Err(err).Msg("postgres connection")
		return
	}

	opts := make([]sequencer.ApiOption, 0)
	if cfg.Indexer.Sequencer.Rps > 0 {
		opts = append(opts, sequencer.WithRateLimit(cfg.Indexer.Sequencer.Rps))
	}
	api := sequencer.NewAPI(cfg.Indexer.Sequencer.Gateway, cfg.Indexer.Sequencer.FeederGateway, opts...)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-signals:
				cancel()
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	testers := []Tester{
		NewJsonSchemaTester(postgres),
		NewTokenBalanceTester(postgres, api),
	}

	for i := range testers {
		if err := testers[i].Test(ctx); err != nil {
			log.Panic().Err(err).Msg(testers[i].String())
			return
		}
	}

	cancel()

	for i := range testers {
		if err := testers[i].Close(); err != nil {
			log.Panic().Err(err).Msg(testers[i].String())
			return
		}
	}

	close(signals)
}
