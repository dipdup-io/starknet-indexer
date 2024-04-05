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
	query := invoke.DB().NewSelect().Model((*storage.Invoke)(nil))
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

	err = invoke.DB().NewSelect().TableExpr("(?) as invoke", query).
		ColumnExpr("invoke.*").
		ColumnExpr("contract.id as contract__id, contract.class_id as contract__class_id, contract.hash as contract__hash, contract.height as contract__height").
		Join("left join address as contract on contract.id = invoke.contract_id").
		Scan(ctx, &result)
	return result, err
}
