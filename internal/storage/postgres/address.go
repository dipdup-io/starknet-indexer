package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10"
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
	err = a.DB().ModelContext(ctx, &address).
		Where("hash = ?", hash).
		Select(&address)
	return
}

// GetAddresses -
func (a *Address) GetAddresses(ctx context.Context, ids ...uint64) (address []storage.Address, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	err = a.DB().ModelContext(ctx, (*storage.Address)(nil)).
		Where("id IN (?)", pg.In(ids)).
		Select(&address)
	return
}

// GetIdsByHash -
func (a *Address) GetIdsByHash(ctx context.Context, hash [][]byte) (ids []uint64, err error) {
	if len(hash) == 0 {
		return
	}

	err = a.DB().ModelContext(ctx, (*storage.Address)(nil)).
		Column("id").
		Where("hash in (?)", pg.In(hash)).
		Select(&ids)
	return
}
