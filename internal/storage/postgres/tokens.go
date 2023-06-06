package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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

// TODO: implement Filterable
