package block

import (
	"context"
	"encoding/json"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

// GetLastBlock -
func GetLastBlock(storage postgres.Storage, ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	block, err := storage.Blocks.Last(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching last block")
	}
	jsonBlock, err := json.Marshal(block)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling json")
	}
	return mcp.NewToolResultText(string(jsonBlock)), nil
}
