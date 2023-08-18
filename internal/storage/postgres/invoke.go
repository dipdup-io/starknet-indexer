package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Invoke -
type Invoke struct {
	*postgres.Table[*storage.Invoke]
}

// NewInvoke -
func NewInvoke(db *database.Bun) *Invoke {
	return &Invoke{
		Table: postgres.NewTable[*storage.Invoke](db),
	}
}

// Filter -
func (invoke *Invoke) Filter(ctx context.Context, fltr []storage.InvokeFilter, opts ...storage.FilterOption) (result []storage.Invoke, err error) {
	query := invoke.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "invoke.id", fltr[i].ID)
				q = integerFilter(q, "invoke.height", fltr[i].Height)
				q = timeFilter(q, "invoke.time", fltr[i].Time)
				q = enumFilter(q, "invoke.status", fltr[i].Status)
				q = enumFilter(q, "invoke.version", fltr[i].Version)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = equalityFilter(q, "invoke.selector", fltr[i].Selector)
				q = stringFilter(q, "invoke.entrypoint", fltr[i].Entrypoint)
				q = jsonFilter(q, "invoke.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "invoke", opts...)
	query.Relation("Contract")

	err = query.Scan(ctx)
	return result, err
}
