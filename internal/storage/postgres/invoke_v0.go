package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// InvokeV0 -
type InvokeV0 struct {
	*postgres.Table[*storage.InvokeV0]
}

// NewInvokeV0 -
func NewInvokeV0(db *database.PgGo) *InvokeV0 {
	return &InvokeV0{
		Table: postgres.NewTable[*storage.InvokeV0](db),
	}
}
