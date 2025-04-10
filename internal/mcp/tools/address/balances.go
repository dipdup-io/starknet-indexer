package address

import (
	"context"
	"encoding/json"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

// GetAddressBalances -
func GetAddressBalances(storage postgres.Storage, ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stringAddress, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, errors.Errorf("address must be a string")
	}

	hashBytes, err := models.HexToBytes(stringAddress)
	if err != nil {
		return nil, err
	}

	address, err := storage.Address.GetByHash(ctx, hashBytes)

	if err != nil {
		return nil, errors.Wrapf(err, "error during executing filter on addresses")
	}

	jsonInvokes, err := json.MarshalIndent(address, "", "  ")
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling invoke")
	}
	return mcp.NewToolResultText(string(jsonInvokes)), nil

}
