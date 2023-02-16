package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Deploy -
type Deploy struct {
	*postgres.Table[*storage.Deploy]
}

// NewDeploy -
func NewDeploy(db *database.PgGo) *Deploy {
	return &Deploy{
		Table: postgres.NewTable[*storage.Deploy](db),
	}
}
