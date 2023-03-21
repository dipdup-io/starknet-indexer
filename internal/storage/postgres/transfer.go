package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Transfer -
type Transfer struct {
	*postgres.Table[*storage.Transfer]
}

// NewTransfer -
func NewTransfer(db *database.PgGo) *Transfer {
	return &Transfer{
		Table: postgres.NewTable[*storage.Transfer](db),
	}
}
