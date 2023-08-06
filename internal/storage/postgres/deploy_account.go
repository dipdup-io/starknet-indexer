package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
)

// DeployAccount -
type DeployAccount struct {
	*postgres.Table[*storage.DeployAccount]
}

// NewDeployAccount -
func NewDeployAccount(db *database.PgGo) *DeployAccount {
	return &DeployAccount{
		Table: postgres.NewTable[*storage.DeployAccount](db),
	}
}

// Filter -
func (d *DeployAccount) Filter(ctx context.Context, fltr []storage.DeployAccountFilter, opts ...storage.FilterOption) ([]storage.DeployAccount, error) {
	query := d.DB().ModelContext(ctx, (*storage.DeployAccount)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "deploy_account.id", fltr[i].ID)
				q = integerFilter(q, "deploy_account.height", fltr[i].Height)
				q = timeFilter(q, "deploy_account.time", fltr[i].Time)
				q = enumFilter(q, "deploy_account.status", fltr[i].Status)
				q = addressFilter(q, "deploy_account.class_id", fltr[i].Class, "Class")
				q = jsonFilter(q, "deploy_account.parsed_calldata", fltr[i].ParsedCalldata)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "deploy_account", opts...)
	query.Relation("Contract").Relation("Class")

	var result []storage.DeployAccount
	err := query.Select(&result)
	return result, err
}
