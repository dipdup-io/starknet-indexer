package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Class -
type Class struct {
	*postgres.Table[*storage.Class]
}

// NewClass -
func NewClass(db *database.PgGo) *Class {
	return &Class{
		Table: postgres.NewTable[*storage.Class](db),
	}
}

// GetByHash -
func (c *Class) GetByHash(ctx context.Context, hash []byte) (class storage.Class, err error) {
	err = c.DB().Model(&class).
		Where("hash = ?", hash).
		Select(&class)
	return
}
