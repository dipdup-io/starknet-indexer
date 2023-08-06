package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
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
func (d *Declare) Filter(ctx context.Context, fltr []storage.DeclareFilter, opts ...storage.FilterOption) ([]storage.Declare, error) {
	query := d.DB().ModelContext(ctx, (*storage.Declare)(nil))

	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "declare.id", fltr[i].ID)
				q = integerFilter(q, "declare.height", fltr[i].Height)
				q = timeFilter(q, "declare.time", fltr[i].Time)
				q = enumFilter(q, "declare.status", fltr[i].Status)
				return enumFilter(q, "declare.version", fltr[i].Version), nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "declare", opts...)
	query.Relation("Contract").Relation("Sender").Relation("Class")

	var declares []storage.Declare
	err := query.Select(&declares)
	return declares, err
}
