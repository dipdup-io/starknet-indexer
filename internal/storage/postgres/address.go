package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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

// Filter -
func (a *Address) Filter(ctx context.Context, fltr []storage.AddressFilter, opts ...storage.FilterOption) ([]storage.Address, error) {
	query := a.DB().ModelContext(ctx, (*storage.Address)(nil))

	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "address.id", fltr[i].ID)
				q = integerFilter(q, "address.height", fltr[i].Height)

				if fltr[i].OnlyStarknet {
					q = q.Where("address.class_id is not null")
				}
				return q, nil
			})
		}
		return q1, nil
	})

	query = optionsFilter(query, "address", opts...)

	var result []storage.Address
	err := query.Select(&result)
	return result, err
}
