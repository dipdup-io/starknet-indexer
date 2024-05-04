package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Event -
type Event struct {
	*postgres.Table[*storage.Event]
}

// NewEvent -
func NewEvent(db *database.Bun) *Event {
	return &Event{
		Table: postgres.NewTable[*storage.Event](db),
	}
}

// Filter -
func (event *Event) Filter(ctx context.Context, fltr []storage.EventFilter, opts ...storage.FilterOption) (result []storage.Event, err error) {
	query := event.DB().NewSelect().Model((*storage.Event)(nil))
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "id", fltr[i].ID)
				q = integerFilter(q, "height", fltr[i].Height)
				q = timeFilter(q, "time", fltr[i].Time)
				q = idFilter(q, "contract_id", fltr[i].Contract)
				q = idFilter(q, "from_id", fltr[i].From)
				q = stringFilter(q, "name", fltr[i].Name)
				q = jsonFilter(q, "parsed_data", fltr[i].ParsedData)
				return q
			})
		}
		return q1
	})

	query = optionsFilter(query, "event", opts...)

	var opt storage.FilterOptions
	for i := range opts {
		opts[i](&opt)
	}

	q := event.DB().NewSelect().
		TableExpr("(?) as event", query).
		ColumnExpr("event.*").
		ColumnExpr("contract.id as contract__id, contract.class_id as contract__class_id, contract.height as contract__height, contract.hash as contract__hash").
		ColumnExpr("from_addr.id as from__id, from_addr.class_id as from__class_id, from_addr.height as from__height, from_addr.hash as from__hash").
		Join("left join address as contract on contract.id = event.contract_id").
		Join("left join address as from_addr on from_addr.id = event.from_id")
	q = addSort(q, opt.SortField, opt.SortOrder)

	err = q.Scan(ctx, &result)
	return
}
