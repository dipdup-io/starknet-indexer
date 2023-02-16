package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// L1Handler -
type L1Handler struct {
	*postgres.Table[*storage.L1Handler]
}

// NewL1Handler -
func NewL1Handler(db *database.PgGo) *L1Handler {
	return &L1Handler{
		Table: postgres.NewTable[*storage.L1Handler](db),
	}
}
