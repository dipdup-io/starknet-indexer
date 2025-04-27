package block

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

// GetLastBlock -
func GetLastBlock(ctx context.Context, storage postgres.Storage, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	block, err := storage.Blocks.Last(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching last block")
	}
	jsonBlock, err := marshalBlock(block)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling json")
	}
	return mcp.NewToolResultText(string(jsonBlock)), nil
}
