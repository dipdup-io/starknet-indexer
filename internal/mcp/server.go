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
}

type Server struct {
	Server   *server.MCPServer
	DbConfig config.Database
	storage  postgres.Storage
}

func NewMCPServer(ctx context.Context, cfg Config) (*Server, error) {
	mcpServer := server.NewMCPServer(
		"example-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)
	postgresStorage, err := postgres.Create(ctx, cfg.Database)
	if err != nil {
		return nil, errors.Wrapf(err, "postgres connection")
	}

	s := &Server{
		Server:   mcpServer,
		DbConfig: cfg.Database,
		storage:  postgresStorage,
	}
	s.addTools()

	return s, nil
}

func (s *Server) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.Server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
	)
}

func (s *Server) addTools() {
	s.Server.AddTool(
		mcp.NewTool(
			"list_tables",
			mcp.WithDescription("List all tables in the database"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.ListTablesTool(s.DbConfig, ctx, req)
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
			return block.GetBlockByHeight(s.storage, ctx, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool("last_block", mcp.WithDescription("Get last block")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return block.GetLastBlock(s.storage, ctx, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_transaction_by_hash",
			mcp.WithDescription("Get transaction data by hash"),
			mcp.WithString("hash",
				mcp.Required(),
				mcp.Description("Transaction hash"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return txs.GetTxByHash(s.storage, ctx, req)
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
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return address.GetAddressBalances(s.storage, ctx, req)
		},
	)
}
