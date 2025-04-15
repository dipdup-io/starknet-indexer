package address

import (
	"context"
	"encoding/json"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type TokenBalance struct {
	OwnerAddress models.HexBytes `json:"owner_address"`
	ContractHash models.HexBytes `json:"contract_hash"`
	TokenId      decimal.Decimal `json:"token_id"`
	Balance      decimal.Decimal `json:"balance"`
}

// GetAddressBalances -
func GetAddressBalances(storage postgres.Storage, ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	address, ok := request.Params.Arguments["address"].(string)
	if !ok {
		return nil, errors.Errorf("address must be a string")
	}
	contract, ok := request.Params.Arguments["contract"].(string)

	addressHash, err := models.HexToBytes(address)
	if err != nil {
		return nil, err
	}
	contractHash, err := models.HexToBytes(contract)
	if err != nil {
		return nil, err
	}

	tokenBalances, err := storage.TokenBalance.Filter(ctx,
		[]models.TokenBalanceFilter{
			{
				Owner: models.BytesFilter{
					Eq: addressHash,
				},
				Contract: models.BytesFilter{
					Eq: contractHash,
				},
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error filtering token balances")
	}

	resultTokenBalances := make([]*TokenBalance, len(tokenBalances))
	for i := range tokenBalances {
		resultTokenBalances[i] = &TokenBalance{
			OwnerAddress: tokenBalances[i].Owner.Hash,
			ContractHash: tokenBalances[i].Contract.Hash,
			TokenId:      tokenBalances[i].TokenID,
			Balance:      tokenBalances[i].Balance,
		}
	}

	jsonTokenBalances, err := json.MarshalIndent(resultTokenBalances, "", "  ")
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling token balances")
	}
	return mcp.NewToolResultText(string(jsonTokenBalances)), nil
}
