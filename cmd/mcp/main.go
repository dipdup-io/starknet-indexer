package main

import (
	"context"
	"os"

	"github.com/dipdup-io/starknet-indexer/internal/mcp"
	"github.com/dipdup-net/go-lib/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog/log"
)

var (
	rootCmd = &cobra.Command{
		Use:   "MCP",
		Short: "DipDup MCP",
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg mcp.Config
	if err := config.Parse(*configPath, &cfg); err != nil {
		log.Panic().Err(err).Msg("parsing config file")
		return
	}

	mcpServer, err := mcp.NewMCPServer(ctx, cfg)
	if err != nil {
		log.Panic().Err(err)
		return
	}

	sseServer := mcpServer.ServeSSE()
	log.Printf("SSE server listening on :8889")

	if err := sseServer.Start(":8889"); err != nil {
		log.Err(err).Msg("Server error: %v")
	}
}
