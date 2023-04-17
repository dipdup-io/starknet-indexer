package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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
func (msg *Message) Filter(ctx context.Context, fltr storage.MessageFilter, opts ...storage.FilterOption) ([]storage.Message, error) {
	q := msg.DB().ModelContext(ctx, (*storage.Message)(nil))
	q = integerFilter(q, "message.id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = addressFilter(q, "hash", fltr.Contract, "Contract")
	q = addressFilter(q, "hash", fltr.From, "From")
	q = addressFilter(q, "hash", fltr.To, "To")
	q = equalityFilter(q, "selector", fltr.Selector)
	q = optionsFilter(q, opts...)

	var result []storage.Message
	err := q.Select(&result)
	return result, err
}
