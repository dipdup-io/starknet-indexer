package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// TokenBalance -
type TokenBalance struct {
	*postgres.Table[*storage.TokenBalance]
}

// NewTokenBalance -
func NewTokenBalance(db *database.Bun) *TokenBalance {
	return &TokenBalance{
		Table: postgres.NewTable[*storage.TokenBalance](db),
	}
}

// Balances -
func (tb *TokenBalance) Balances(ctx context.Context, contractId uint64, tokenId int64, limit, offset int) (balances []storage.TokenBalance, err error) {
	query := tb.DB().NewSelect().Model(&balances).
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
		Scan(ctx)
	return
}

// Filter -
func (tb *TokenBalance) Filter(ctx context.Context, fltr []storage.TokenBalanceFilter, opts ...storage.FilterOption) (result []storage.TokenBalance, err error) {
	query := tb.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" OR ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = addressFilter(q, "hash", fltr[i].Owner, "Owner")
				q = stringFilter(q, "token_balance.token_id", fltr[i].TokenId)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "token_balance", opts...)
	query.Relation("Contract").Relation("Owner")

	err = query.Scan(ctx)
	return result, err
}
