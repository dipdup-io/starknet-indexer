package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
)

// Message -
type Message struct {
	*postgres.Table[*storage.Message]
}

// NewMessage -
func NewMessage(db *database.PgGo) *Message {
	return &Message{
		Table: postgres.NewTable[*storage.Message](db),
	}
}

// Filter -
func (msg *Message) Filter(ctx context.Context, fltr []storage.MessageFilter, opts ...storage.FilterOption) ([]storage.Message, error) {
	query := msg.DB().ModelContext(ctx, (*storage.Message)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "message.id", fltr[i].ID)
				q = integerFilter(q, "message.height", fltr[i].Height)
				q = timeFilter(q, "message.time", fltr[i].Time)
				q = addressFilter(q, "message.contract_id", fltr[i].Contract, "Contract")
				q = addressFilter(q, "message.from_id", fltr[i].From, "From")
				q = addressFilter(q, "message.to_id", fltr[i].To, "To")
				q = equalityFilter(q, "message.selector", fltr[i].Selector)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "message", opts...)

	var result []storage.Message
	err := query.Select(&result)
	return result, err
}
