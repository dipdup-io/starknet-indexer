package postgres

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Internal -
type Internal struct {
	*postgres.Table[*storage.Internal]
}

// NewInternal -
func NewInternal(db *database.Bun) *Internal {
	return &Internal{
		Table: postgres.NewTable[*storage.Internal](db),
	}
}

// Filter -
func (d *Internal) Filter(ctx context.Context, fltr []storage.InternalFilter, opts ...storage.FilterOption) (result []storage.Internal, err error) {
	query := d.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "internal_tx.id", fltr[i].ID)
				q = integerFilter(q, "internal_tx.height", fltr[i].Height)
				q = timeFilter(q, "internal_tx.time", fltr[i].Time)
				q = enumFilter(q, "internal_tx.status", fltr[i].Status)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = addressFilter(q, "hash", fltr[i].Caller, "Caller")
				q = addressFilter(q, "hash", fltr[i].Class, "Class")
				q = equalityFilter(q, "internal_tx.selector", fltr[i].Selector)
				q = stringFilter(q, "internal_tx.entrypoint", fltr[i].Entrypoint)
				q = enumFilter(q, "internal_tx.entrypoint_type", fltr[i].EntrypointType)
				q = enumFilter(q, "internal_tx.call_type", fltr[i].CallType)
				q = jsonFilter(q, "internal_tx.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "internal_tx", opts...)
	query.Relation("Contract").Relation("Caller").Relation("Class")

	err = query.Scan(ctx)
	return
}

func (i *Internal) GetDeployedContracts(ctx context.Context, deployerBytes []byte) ([]storage.DeployedContract, error) {
	var contracts []storage.DeployedContract
	err := i.DB().NewRaw(`
		SELECT DISTINCT ON (d.hash)
			d.time as deploy_time,
			deployed.hash as contract_address,
			d.hash as tx_hash
		FROM internal_tx it
		JOIN deploy d ON it.deploy_id = d.id
		JOIN address deployer ON it.caller_id = deployer.id
		JOIN address deployed ON d.contract_id = deployed.id
		WHERE it.deploy_id IS NOT NULL
		AND deployer.hash = ?
		ORDER BY d.hash, d.time DESC
	`, deployerBytes).Scan(ctx, &contracts)

	return contracts, err
}
