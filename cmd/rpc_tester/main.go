package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/printer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "indexer",
		Short: "DipDup indexer",
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

	client := grpc.NewClient(*cfg.GRPC)
	print := NewPrinter()

	if err := modules.Connect(client, print, grpc.OutputMessages, printer.InputName); err != nil {
		log.Panic().Err(err).Msg("module connect")
		return
	}

	if err := client.Connect(ctx); err != nil {
		log.Panic().Err(err).Msg("grpc connect")
		return
	}

	client.Start(ctx)
	print.Start(ctx)

	subscriptions := make([]uint64, 0)
	for name, sub := range cfg.GRPC.Subscriptions {
		log.Info().Str("topic", name).Msg("subscribing...")
		req := sub.ToGrpcFilter()
		subId, err := client.Subscribe(ctx, req)
		if err != nil {
			log.Panic().Err(err).Msg("subscribing error")
			return
		}
		subscriptions = append(subscriptions, subId)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	for i := range subscriptions {
		if err := client.Unsubscribe(ctx, subscriptions[i]); err != nil {
			log.Panic().Err(err).Msg("unsubscribing")
		}
	}

	cancel()

	if err := print.Close(); err != nil {
		log.Panic().Err(err).Msg("closing printer")
	}
	if err := client.Close(); err != nil {
		log.Panic().Err(err).Msg("closing grpc server")
	}

	close(signals)
}
