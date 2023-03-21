package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/shopspring/decimal"
)

// TokenBalance -
type TokenBalance struct {
	*postgres.Table[*storage.TokenBalance]
}

// NewTokenBalance -
func NewTokenBalance(db *database.PgGo) *TokenBalance {
	return &TokenBalance{
		Table: postgres.NewTable[*storage.TokenBalance](db),
	}
}

// NegativeBalances -
func (tb *TokenBalance) NegativeBalances(ctx context.Context) (balances []storage.TokenBalance, err error) {
	err = tb.DB().ModelContext(ctx, &balances).
		Where("balance < 0").
		Relation("Contract").
		Relation("Owner").
		Select(&balances)
	return
}

type supply struct {
	Value string `pg:"value"`
}

// TotalSupply -
func (tb *TokenBalance) TotalSupply(ctx context.Context, contractId, tokenId uint64) (decimal.Decimal, error) {
	var supply supply
	if err := tb.DB().ModelContext(ctx, (*storage.TokenBalance)(nil)).
		ColumnExpr("sum(balance)::text as value").
		Where("balance > 0").
		Where("contract_id = ?", contractId).
		Where("token_id = ?", tokenId).
		Select(&supply); err != nil {
		return decimal.Zero, err
	}

	if supply.Value == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(supply.Value)
}

// Owner -
func (tb *TokenBalance) Owner(ctx context.Context, cotractId uint64, tokenId decimal.Decimal) (tokenBalance storage.TokenBalance, err error) {
	err = tb.DB().ModelContext(ctx, &tokenBalance).
		Where("contract_id = ?", cotractId).
		Where("token_id = ?", tokenId).
		Limit(1).
		Relation("Owner").
		Relation("Contract").
		Select(&tokenBalance)
	return
}
