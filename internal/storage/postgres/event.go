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
	query := event.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "event.id", fltr[i].ID)
				q = integerFilter(q, "event.height", fltr[i].Height)
				q = timeFilter(q, "event.time", fltr[i].Time)
				q = idFilter(q, "event.contract_id", fltr[i].Contract, "Contract")
				q = idFilter(q, "event.from_id", fltr[i].From, "From")
				q = stringFilter(q, "event.name", fltr[i].Name)
				q = jsonFilter(q, "event.parsed_data", fltr[i].ParsedData)
				return q
			})
		}
		return q1
	})

	query = optionsFilter(query, "event", opts...)
	query.Relation("Contract").Relation("From")

	err = query.Scan(ctx)
	return
}
