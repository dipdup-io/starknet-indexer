package block

import (
	"context"
	"encoding/json"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"strconv"
)

// GetBlockByHeight -
func GetBlockByHeight(storage postgres.Storage, ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	requestHeight, ok := request.Params.Arguments["height"].(string)
	if !ok {
		return nil, errors.Errorf("height must be a string")
	}
	height, err := strconv.ParseUint(requestHeight, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "can't convert tool request param height to uint64")
	}

	block, err := storage.Blocks.ByHeight(ctx, height)
	if err != nil {
		return nil, errors.Wrapf(err, "error fetching block with height %d", height)
	}

	jsonBlock, err := json.Marshal(block)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling json")
	}
	return mcp.NewToolResultText(string(jsonBlock)), nil
}
