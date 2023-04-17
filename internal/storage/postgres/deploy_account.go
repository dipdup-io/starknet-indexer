package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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
func (d *DeployAccount) Filter(ctx context.Context, fltr storage.DeployAccountFilter, opts ...storage.FilterOption) ([]storage.DeployAccount, error) {
	q := d.DB().ModelContext(ctx, (*storage.DeployAccount)(nil))
	q = integerFilter(q, "deploy_account.id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = enumFilter(q, "status", fltr.Status)
	q = addressFilter(q, "hash", fltr.Class, "Class")
	q = jsonFilter(q, "parsed_calldata", fltr.ParsedCalldata)
	q = optionsFilter(q, opts...)

	var result []storage.DeployAccount
	err := q.Select(&result)
	return result, err
}
