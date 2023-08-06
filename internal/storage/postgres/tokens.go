package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
)

// Token -
type Token struct {
	*postgres.Table[*storage.Token]
}

// NewToken -
func NewToken(db *database.PgGo) *Token {
	return &Token{
		Table: postgres.NewTable[*storage.Token](db),
	}
}

// ListByType -
func (tokens *Token) ListByType(ctx context.Context, typ storage.TokenType, limit uint64, offset uint64, order sdk.SortOrder) ([]storage.Token, error) {
	query := tokens.DB().ModelContext(ctx, (*storage.Token)(nil)).
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

	var result []storage.Token
	err := query.Select(&result)
	return result, err
}

// Filter -
func (token *Token) Filter(ctx context.Context, fltr []storage.TokenFilter, opts ...storage.FilterOption) ([]storage.Token, error) {
	query := token.DB().ModelContext(ctx, (*storage.Token)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "token.id", fltr[i].ID)
				q = addressFilter(q, "token.contract_id", fltr[i].Contract, "Contract")
				q = stringFilter(q, "token.token_id", fltr[i].TokenId)
				q = enumStringFilter(q, "token.type", fltr[i].Type)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "token", opts...)
	query = query.Relation("Contract")

	var result []storage.Token
	err := query.Select(&result)
	return result, err
}

// Find -
func (token *Token) Find(ctx context.Context, contractId uint64, tokenId string) (t storage.Token, err error) {
	err = token.DB().ModelContext(ctx, &t).
		Where("contract_id = ?", contractId).
		Where("token_id = ?", tokenId).
		Limit(1).
		Select(&t)
	return
}
