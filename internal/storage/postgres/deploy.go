package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Deploy -
type Deploy struct {
	*postgres.Table[*storage.Deploy]
}

// NewDeploy -
func NewDeploy(db *database.Bun) *Deploy {
	return &Deploy{
		Table: postgres.NewTable[*storage.Deploy](db),
	}
}

// Filter -
func (d *Deploy) Filter(ctx context.Context, fltr []storage.DeployFilter, opts ...storage.FilterOption) (result []storage.Deploy, err error) {
	query := d.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "deploy.id", fltr[i].ID)
				q = integerFilter(q, "deploy.height", fltr[i].Height)
				q = timeFilter(q, "deploy.time", fltr[i].Time)
				q = enumFilter(q, "deploy.status", fltr[i].Status)
				q = addressFilter(q, "hash", fltr[i].Class, "Class")
				q = jsonFilter(q, "deploy.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "deploy", opts...)
	query.Relation("Contract").Relation("Class")

	err = query.Scan(ctx)
	return
}
