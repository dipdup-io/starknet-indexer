package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/spf13/cobra"

	"net/http"
	_ "net/http/pprof"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	rootCmd = &cobra.Command{
		Use:   "indexer",
		Short: "DipDup indexer",
	}
)

func main() {
	go func() {
		log.Print(http.ListenAndServe("localhost:6060", nil))
	}()
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

	if err := starknet.LoadBridgedTokens(cfg.Indexer.BridgedTokensFile); err != nil {
		log.Panic().Err(err).Msg("loading bridged tokens")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	postgres, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		log.Panic().Err(err).Msg("postgres connection")
		return
	}

	indexerModule := indexer.New(cfg.Indexer, postgres)

	grpcModule, err := grpc.NewServer(
		cfg.GRPC,
	)
	if err != nil {
		log.Panic().Err(err).Msg("creating grpc module")
		cancel()
		return
	}

	if err := modules.Connect(indexerModule, grpcModule, indexer.OutputBlocks, grpc.InputBlocks); err != nil {
		log.Panic().Err(err).Msg("creating modules connection")
		cancel()
		return
	}

	grpcModule.Start(ctx)
	indexerModule.Start(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	cancel()

	if err := indexerModule.Close(); err != nil {
		log.Panic().Err(err).Msg("closing indexer")
	}
	if err := grpcModule.Close(); err != nil {
		log.Panic().Err(err).Msg("closing grpc server")
	}

	close(signals)
}
