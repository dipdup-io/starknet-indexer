package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Token -
type Token struct {
	*postgres.Table[*storage.Token]
}

// NewToken -
func NewToken(db *database.Bun) *Token {
	return &Token{
		Table: postgres.NewTable[*storage.Token](db),
	}
}

// ListByType -
func (tokens *Token) ListByType(ctx context.Context, typ storage.TokenType, limit uint64, offset uint64, order sdk.SortOrder) (result []storage.Token, err error) {
	query := tokens.DB().NewSelect().Model(&result).
		Where("type = ?", typ).
		Offset(int(offset))
	if limit == 0 {
		query.Limit(10)
	} else {
		query.Limit(int(limit))
	}

	switch order {
	case sdk.SortOrderAsc:
		query.Order("id asc")
	case sdk.SortOrderDesc:
		query.Order("id desc")
	default:
		query.Order("id asc")
	}

	err = query.Scan(ctx)
	return result, err
}

// Filter -
func (token *Token) Filter(ctx context.Context, fltr []storage.TokenFilter, opts ...storage.FilterOption) (result []storage.Token, err error) {
	query := token.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "token.id", fltr[i].ID)
				q = idFilter(q, "contract_id", fltr[i].Contract)
				q = stringFilter(q, "token.token_id", fltr[i].TokenId)
				q = enumStringFilter(q, "token.type", fltr[i].Type)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "token", opts...)

	var opt storage.FilterOptions
	for i := range opts {
		opts[i](&opt)
	}

	q := token.DB().NewSelect().
		TableExpr("(?) as token", query).
		ColumnExpr("token.*").
		ColumnExpr("contract.id as contract__id, contract.class_id as contract__class_id, contract.height as contract__height, contract.hash as contract__hash").
		Join("left join address as contract on contract.id = token.contract_id")

	q = addSort(q, opt.SortField, opt.SortOrder)

	err = q.Scan(ctx, &result)
	return result, err
}

// Find -
func (token *Token) Find(ctx context.Context, contractId uint64, tokenId string) (t storage.Token, err error) {
	err = token.DB().NewSelect().Model(&t).
		Where("contract_id = ?", contractId).
		Where("token_id = ?", tokenId).
		Limit(1).
		Scan(ctx)
	return
}
