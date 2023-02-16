package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// DeployAccount -
type DeployAccount struct {
	*postgres.Table[*storage.DeployAccount]
}

// NewDeployAccount -
func NewDeployAccount(db *database.PgGo) *DeployAccount {
	return &DeployAccount{
		Table: postgres.NewTable[*storage.DeployAccount](db),
	}
}
