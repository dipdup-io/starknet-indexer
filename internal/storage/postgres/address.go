package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// Address -
type Address struct {
	*postgres.Table[*storage.Address]
}

// NewAddress -
func NewAddress(db *database.PgGo) *Address {
	return &Address{
		Table: postgres.NewTable[*storage.Address](db),
	}
}

// GetByHash -
func (a *Address) GetByHash(ctx context.Context, hash []byte) (address storage.Address, err error) {
	err = a.DB().Model(&address).
		Where("hash = ?", hash).
		Select(&address)
	return
}
