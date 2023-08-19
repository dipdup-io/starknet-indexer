package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Declare -
type Declare struct {
	*postgres.Table[*storage.Declare]
}

// NewDeclare -
func NewDeclare(db *database.Bun) *Declare {
	return &Declare{
		Table: postgres.NewTable[*storage.Declare](db),
	}
}

// Filter -
func (d *Declare) Filter(ctx context.Context, fltr []storage.DeclareFilter, opts ...storage.FilterOption) (result []storage.Declare, err error) {
	query := d.DB().NewSelect().Model(&result)

	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "declare.id", fltr[i].ID)
				q = integerFilter(q, "declare.height", fltr[i].Height)
				q = timeFilter(q, "declare.time", fltr[i].Time)
				q = enumFilter(q, "declare.status", fltr[i].Status)
				return enumFilter(q, "declare.version", fltr[i].Version)
			})
		}
		return q1
	})
	query = optionsFilter(query, "declare", opts...)
	query.Relation("Contract").Relation("Sender").Relation("Class")

	err = query.Scan(ctx)
	return
}
