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

// Filter -
func (d *Declare) Filter(ctx context.Context, fltr storage.DeclareFilter, opts ...storage.FilterOption) ([]storage.Declare, error) {
	q := d.DB().ModelContext(ctx, (*storage.Declare)(nil))
	q = integerFilter(q, "id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = enumFilter(q, "status", fltr.Status)
	q = enumFilter(q, "version", fltr.Version)
	q = optionsFilter(q, opts...)

	var declares []storage.Declare
	err := q.Select(&declares)
	return declares, err
}
