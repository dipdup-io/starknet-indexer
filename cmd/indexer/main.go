package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/hasura"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/spf13/cobra"

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
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.Logger.With().Caller().Logger()

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

	views, err := createViews(postgres)
	if err != nil {
		log.Panic().Err(err).Msg("create views")
		return
	}

	if cfg.Hasura != nil {
		if err := hasura.Create(ctx, hasura.GenerateArgs{
			Config:         cfg.Hasura,
			DatabaseConfig: cfg.Database,
			Models:         storage.ModelsAny,
			Views:          views,
		}); err != nil {
			log.Panic().Err(err).Msg("hasura initialization")
			return
		}
	}

	indexerModule, err := indexer.New(cfg.Indexer, postgres, cfg.DataSources)
	if err != nil {
		log.Panic().Err(err).Msg("creating indexer module")
		cancel()
		return
	}

	grpcModule, err := grpc.NewServer(
		cfg.GRPC, postgres,
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

	notifyCtx, notifyCancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer notifyCancel()

	<-notifyCtx.Done()
	cancel()

	if err := indexerModule.Close(); err != nil {
		log.Panic().Err(err).Msg("closing indexer")
	}
	if err := grpcModule.Close(); err != nil {
		log.Panic().Err(err).Msg("closing grpc server")
	}
}
