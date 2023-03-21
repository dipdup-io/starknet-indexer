package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Fee -
type Fee struct {
	*postgres.Table[*storage.Fee]
}

// NewFee -
func NewFee(db *database.PgGo) *Fee {
	return &Fee{
		Table: postgres.NewTable[*storage.Fee](db),
	}
}
