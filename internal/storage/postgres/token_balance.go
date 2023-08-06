package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
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
func (tb *TokenBalance) Owner(ctx context.Context, contractId uint64, tokenId decimal.Decimal) (tokenBalance storage.TokenBalance, err error) {
	err = tb.DB().ModelContext(ctx, &tokenBalance).
		Where("contract_id = ?", contractId).
		Where("token_id = ?", tokenId).
		Where("balance > 0").
		Limit(1).
		Relation("Owner").
		Relation("Contract").
		Select(&tokenBalance)
	return
}

// Balances -
func (tb *TokenBalance) Balances(ctx context.Context, contractId uint64, tokenId int64, limit, offset int) (balances []storage.TokenBalance, err error) {
	query := tb.DB().ModelContext(ctx, &balances).
		Where("contract_id = ?", contractId)

	if tokenId >= 0 {
		query.Where("token_id = ?", tokenId)
	}
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	err = query.
		Relation("Owner").
		Relation("Contract").
		Select(&balances)
	return
}

// Filter -
func (tb *TokenBalance) Filter(ctx context.Context, fltr []storage.TokenBalanceFilter, opts ...storage.FilterOption) ([]storage.TokenBalance, error) {
	query := tb.DB().ModelContext(ctx, (*storage.TokenBalance)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = addressFilter(q, "token_balance.contract_id", fltr[i].Contract, "Contract")
				q = addressFilter(q, "token_balance.owner_id", fltr[i].Owner, "Owner")
				q = stringFilter(q, "token_balance.token_id", fltr[i].TokenId)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "token_balance", opts...)
	query.Relation("Contract").Relation("Owner")

	var result []storage.TokenBalance
	err := query.Select(&result)
	return result, err
}
