package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Message -
type Message struct {
	*postgres.Table[*storage.Message]
}

// NewMessage -
func NewMessage(db *database.Bun) *Message {
	return &Message{
		Table: postgres.NewTable[*storage.Message](db),
	}
}

// Filter -
func (msg *Message) Filter(ctx context.Context, fltr []storage.MessageFilter, opts ...storage.FilterOption) (result []storage.Message, err error) {
	query := msg.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "message.id", fltr[i].ID)
				q = integerFilter(q, "message.height", fltr[i].Height)
				q = timeFilter(q, "message.time", fltr[i].Time)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = addressFilter(q, "hash", fltr[i].From, "From")
				q = addressFilter(q, "hash", fltr[i].To, "To")
				q = equalityFilter(q, "hash", fltr[i].Selector)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "message", opts...)
	query.Relation("Contract").Relation("From").Relation("To")

	err = query.Scan(ctx)
	return result, err
}
