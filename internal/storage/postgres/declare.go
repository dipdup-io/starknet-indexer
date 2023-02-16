package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Declare -
type Declare struct {
	*postgres.Table[*storage.Declare]
}

// NewDeclare -
func NewDeclare(db *database.PgGo) *Declare {
	return &Declare{
		Table: postgres.NewTable[*storage.Declare](db),
	}
}
