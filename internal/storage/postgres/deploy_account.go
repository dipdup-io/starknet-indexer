package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// DeployAccount -
type DeployAccount struct {
	*postgres.Table[*storage.DeployAccount]
}

// NewDeployAccount -
func NewDeployAccount(db *database.Bun) *DeployAccount {
	return &DeployAccount{
		Table: postgres.NewTable[*storage.DeployAccount](db),
	}
}

// Filter -
func (d *DeployAccount) Filter(ctx context.Context, fltr []storage.DeployAccountFilter, opts ...storage.FilterOption) (result []storage.DeployAccount, err error) {
	query := d.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "deploy_account.id", fltr[i].ID)
				q = integerFilter(q, "deploy_account.height", fltr[i].Height)
				q = timeFilter(q, "deploy_account.time", fltr[i].Time)
				q = enumFilter(q, "deploy_account.status", fltr[i].Status)
				q = addressFilter(q, "hash", fltr[i].Class, "Class")
				q = jsonFilter(q, "deploy_account.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "deploy_account", opts...)
	query.Relation("Contract").Relation("Class")

	err = query.Scan(ctx)
	return
}

func (d *DeployAccount) HashByHeight(ctx context.Context, height uint64) (hash []byte, err error) {
	err = d.DB().NewSelect().
		Model((*storage.DeployAccount)(nil)).
		Column("hash").
		Where("height = ?", height).
		Limit(1).
		Scan(ctx, &hash)
	return
}
