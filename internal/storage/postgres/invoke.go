package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Invoke -
type Invoke struct {
	*postgres.Table[*storage.Invoke]
}

// NewInvoke -
func NewInvoke(db *database.PgGo) *Invoke {
	return &Invoke{
		Table: postgres.NewTable[*storage.Invoke](db),
	}
}

// ByHeight -
func (invoke *Invoke) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.Invoke, err error) {
	err = invoke.DB().ModelContext(ctx, (*storage.Invoke)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
