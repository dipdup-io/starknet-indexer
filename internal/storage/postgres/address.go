package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Address -
type Address struct {
	*postgres.Table[*storage.Address]
}

// NewAddress -
func NewAddress(db *database.Bun) *Address {
	return &Address{
		Table: postgres.NewTable[*storage.Address](db),
	}
}

// GetByHash -
func (a *Address) GetByHash(ctx context.Context, hash []byte) (address storage.Address, err error) {
	err = a.DB().NewSelect().Model(&address).
		Where("hash = ?", hash).
		Scan(ctx)
	return
}

// GetAddresses -
func (a *Address) GetAddresses(ctx context.Context, ids ...uint64) (address []storage.Address, err error) {
	if len(ids) == 0 {
		return nil, nil
	}

	err = a.DB().NewSelect().Model(&address).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)
	return
}

// GetIdsByHash -
func (a *Address) GetByHashes(ctx context.Context, hash [][]byte) (addresses []storage.Address, err error) {
	if len(hash) == 0 {
		return
	}

	err = a.DB().NewSelect().Model(&addresses).
		Where("hash in (?)", bun.In(hash)).
		Scan(ctx)
	return
}

// Filter -
func (a *Address) Filter(ctx context.Context, fltr []storage.AddressFilter, opts ...storage.FilterOption) (result []storage.Address, err error) {
	query := a.DB().NewSelect().Model(&result)

	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "address.id", fltr[i].ID)
				q = integerFilter(q, "address.height", fltr[i].Height)

				if fltr[i].OnlyStarknet {
					q = q.Where("address.class_id is not null")
				}
				return q
			})
		}
		return q1
	})

	query = optionsFilter(query, "address", opts...)

	err = query.Scan(ctx)
	return
}
