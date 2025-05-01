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
			mcp.WithDescription("Retrieves all tables in the Starknet database."),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:          "List tables",
				ReadOnlyHint:   true,
				IdempotentHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.ListTablesTool(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_block_by_height",
			mcp.WithDescription("Retrieves Starknet block data for a given height. Block data contains time of block creation, transactions count, cryptographic hashes and etc."),
			mcp.WithString("height",
				mcp.Required(),
				mcp.Description("Numerical identifier assigned to a specific block (e.g. 123456)"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get block by height",
				ReadOnlyHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return block.GetBlockByHeight(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"last_block",
			mcp.WithDescription("Retrieves the last Starknet block. Block data contains time of block creation, transactions count, cryptographic hashes and etc."),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get last block",
				ReadOnlyHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return block.GetLastBlock(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_transaction_by_hash",
			mcp.WithDescription("Retrieves transaction data by its hash. Transaction data contains its status, time, corresponding block height, calldata, entrypoint, called contract, etc."),
			mcp.WithString("hash",
				mcp.Required(),
				mcp.Description("Transaction hash (e.g. 0x4d7746e37e278c17e53b942cdfd9097435c4f9be851b7707f434cc8d1379ca9)"),
			),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get transaction by hash",
				ReadOnlyHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return txs.GetTxByHash(ctx, s.storage, req)
		},
	)
	s.Server.AddTool(
		mcp.NewTool(
			"get_address_balances",
			mcp.WithDescription("Retrieves Starknet address token balances."),
			WithAddress("address", mcp.Required()),
			WithAddress("contract"),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get address balances",
				ReadOnlyHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return address.GetAddressBalances(ctx, s.storage, req)
		},
	)

	s.Server.AddTool(
		mcp.NewTool(
			"get_address_activity",
			mcp.WithDescription("Retrieves complex chain activity address data which containing general address info, last transfers, invocations, deployments and events."),
			WithAddress("address", mcp.Required()),
			mcp.WithString("limit", mcp.Description("Rows limit, default is 20")),
			mcp.WithString("offset", mcp.Description("Offset, default is 0")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        "Get address activity",
				ReadOnlyHint: true,
			}),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return address.GetAddressActivity(ctx, s.storage, req)
		},
	)
}

func WithAddress(name string, opts ...mcp.PropertyOption) mcp.ToolOption {
	opts = append(opts,
		mcp.Description("Starknet address is hash which identify user or contract, e.g. 0x04c53aea8a08338206f37cb9f9a64664193d3f2ad670e84e77230189d60215e5"),
	)
	return mcp.WithString(name, opts...)
}
