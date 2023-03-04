package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Internal -
type Internal struct {
	*postgres.Table[*storage.Internal]
}

// NewInternal -
func NewInternal(db *database.PgGo) *Internal {
	return &Internal{
		Table: postgres.NewTable[*storage.Internal](db),
	}
}
