package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// L1Handler -
type L1Handler struct {
	*postgres.Table[*storage.L1Handler]
}

// NewL1Handler -
func NewL1Handler(db *database.PgGo) *L1Handler {
	return &L1Handler{
		Table: postgres.NewTable[*storage.L1Handler](db),
	}
}

// ByHeight -
func (l1 *L1Handler) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.L1Handler, err error) {
	err = l1.DB().ModelContext(ctx, (*storage.L1Handler)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
