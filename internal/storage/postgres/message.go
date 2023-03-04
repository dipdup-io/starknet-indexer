package postgres

import (
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
