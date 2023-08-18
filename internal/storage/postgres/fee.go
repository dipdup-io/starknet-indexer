package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Fee -
type Fee struct {
	*postgres.Table[*storage.Fee]
}

// NewFee -
func NewFee(db *database.Bun) *Fee {
	return &Fee{
		Table: postgres.NewTable[*storage.Fee](db),
	}
}

// Filter -
func (fee *Fee) Filter(ctx context.Context, fltr []storage.FeeFilter, opts ...storage.FilterOption) (result []storage.Fee, err error) {
	query := fee.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "fee.id", fltr[i].ID)
				q = integerFilter(q, "fee.height", fltr[i].Height)
				q = timeFilter(q, "fee.time", fltr[i].Time)
				q = enumFilter(q, "fee.status", fltr[i].Status)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = addressFilter(q, "hash", fltr[i].Caller, "Caller")
				q = addressFilter(q, "hash", fltr[i].Class, "Class")
				q = equalityFilter(q, "fee.selector", fltr[i].Selector)
				q = stringFilter(q, "fee.entrypoint", fltr[i].Entrypoint)
				q = enumFilter(q, "fee.entrypoint_type", fltr[i].EntrypointType)
				q = enumFilter(q, "fee.call_type", fltr[i].CallType)
				q = jsonFilter(q, "fee.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "fee", opts...)
	query.Relation("Contract").Relation("Caller").Relation("Class")

	err = query.Scan(ctx)
	return
}
