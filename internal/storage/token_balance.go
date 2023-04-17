package storage

import (
	"context"
	"fmt"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// ITokenBalance -
type ITokenBalance interface {
	storage.Table[*TokenBalance]
	Filterable[TokenBalance, TokenBalanceFilter]

	NegativeBalances(ctx context.Context) ([]TokenBalance, error)
	TotalSupply(ctx context.Context, contractId, tokenId uint64) (decimal.Decimal, error)
	Owner(ctx context.Context, cotractId uint64, tokenId decimal.Decimal) (TokenBalance, error)
	Balances(ctx context.Context, contractId uint64, tokenId int64, limit, offset int) ([]TokenBalance, error)
}

// TokenBalanceFilter -
type TokenBalanceFilter struct {
	Owner    BytesFilter
	Contract BytesFilter
	TokenId  StringFilter
}

// TokenBalance -
type TokenBalance struct {
	// nolint
	tableName struct{} `pg:"token_balance"`

	OwnerID    uint64          `pg:",pk"`
	ContractID uint64          `pg:",pk"`
	TokenID    decimal.Decimal `pg:",pk,type:numeric,use_zero"`
	Balance    decimal.Decimal `pg:",type:numeric,use_zero"`

	Contract Address `pg:"rel:has-one"`
	Owner    Address `pg:"rel:has-one"`
}

// TableName -
func (TokenBalance) TableName() string {
	return "token_balance"
}

// String -
func (tb TokenBalance) String() string {
	return fmt.Sprintf("%d_%d_%s", tb.ContractID, tb.OwnerID, tb.TokenID.String())
}
