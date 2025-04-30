package mcp

import (
	"context"
	"fmt"

	"github.com/dipdup-io/starknet-indexer/internal/mcp/tools"
	"github.com/dipdup-io/starknet-indexer/internal/mcp/tools/address"
	"github.com/dipdup-io/starknet-indexer/internal/mcp/tools/block"
	"github.com/dipdup-io/starknet-indexer/internal/mcp/tools/txs"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-net/go-lib/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pkg/errors"
)

type Config struct {
	config.Config `yaml:",inline"`

	Mcp *ServerConfig `yaml:"mcp"`
}

type ServerConfig struct {
	Bind string `validate:"required,hostname_port" yaml:"bind"`
}

type Server struct {
	Server  *server.MCPServer
	storage postgres.Storage
	bind    string
}

func NewMCPServer(ctx context.Context, cfg Config) (*Server, error) {
	if cfg.Mcp == nil {
		return nil, errors.New("config 'mcp' section is absent")
	}
	mcpServer := server.NewMCPServer(
		"starknet-mcp-server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	postgresStorage, err := postgres.Create(ctx, cfg.Database, false)
	if err != nil {
		return nil, errors.Wrapf(err, "postgres connection")
	}

	s := &Server{
		Server:  mcpServer,
		storage: postgresStorage,
		bind:    cfg.Mcp.Bind,
	}
	s.addTools()

	return s, nil
}

func (s *Server) ServeSSE() *server.SSEServer {
	return server.NewSSEServer(s.Server,
		server.WithBaseURL(fmt.Sprintf("http://%s", s.bind)),
		server.WithUseFullURLForMessageEndpoint(false),
	)
}

func (s *Server) addTools() {
	s.Server.AddTool(
		mcp.NewTool(
			"list_tables",
			mcp.WithDescription("List all tables in the database"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.ListTablesTool(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_block_by_height",
			mcp.WithDescription("Get block data by height"),
			mcp.WithString("height",
				mcp.Required(),
				mcp.Description("Block height"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return block.GetBlockByHeight(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool("last_block", mcp.WithDescription("Get last block")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return block.GetLastBlock(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_transaction_by_hash",
			mcp.WithDescription("Get transaction by its hash"),
			mcp.WithString("hash",
				mcp.Required(),
				mcp.Description("Transaction hash"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return txs.GetTxByHash(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_address_balances",
			mcp.WithDescription("Get address token balances"),
			mcp.WithString("address",
				mcp.Required(),
				mcp.Description("Starknet address"),
			),
			mcp.WithString("contract", mcp.Description("Starknet address")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return address.GetAddressBalances(ctx, s.storage, req)
		},
	)

	s.Server.AddTool(
		mcp.NewTool(
			"get_address_activity",
			mcp.WithDescription("Complex chain activity address data and stats"),
			mcp.WithString("address",
				mcp.Required(),
				mcp.Description("Starknet address"),
			),
			mcp.WithString("limit", mcp.Description("Rows limit, default is 20")),
			mcp.WithString("offset", mcp.Description("Offset, default is 0")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return address.GetAddressActivity(ctx, s.storage, req)
		},
	)
}
