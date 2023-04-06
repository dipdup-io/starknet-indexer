package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Declare -
type Declare struct {
	*postgres.Table[*storage.Declare]
}

// NewDeclare -
func NewDeclare(db *database.PgGo) *Declare {
	return &Declare{
		Table: postgres.NewTable[*storage.Declare](db),
	}
}

// ByHeight -
func (d *Declare) ByHeight(ctx context.Context, height, limit, offset uint64) (response []storage.Declare, err error) {
	err = d.DB().ModelContext(ctx, (*storage.Declare)(nil)).
		Where("height = ?", height).
		Limit(int(limit)).
		Offset(int(offset)).
		Select(&response)
	return
}
