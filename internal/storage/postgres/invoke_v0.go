package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// InvokeV0 -
type InvokeV0 struct {
	*postgres.Table[*storage.InvokeV0]
}

// NewInvokeV0 -
func NewInvokeV0(db *database.PgGo) *InvokeV0 {
	return &InvokeV0{
		Table: postgres.NewTable[*storage.InvokeV0](db),
	}
}

// ByHeight -
func (invoke *InvokeV0) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.InvokeV0, err error) {
	err = invoke.DB().ModelContext(ctx, (*storage.InvokeV0)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
