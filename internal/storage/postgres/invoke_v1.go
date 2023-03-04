package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// InvokeV1 -
type InvokeV1 struct {
	*postgres.Table[*storage.InvokeV1]
}

// NewInvokeV1 -
func NewInvokeV1(db *database.PgGo) *InvokeV1 {
	return &InvokeV1{
		Table: postgres.NewTable[*storage.InvokeV1](db),
	}
}

// ByHeight -
func (invoke *InvokeV1) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.InvokeV1, err error) {
	err = invoke.DB().ModelContext(ctx, (*storage.InvokeV1)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
