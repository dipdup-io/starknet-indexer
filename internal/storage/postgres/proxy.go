package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
)

// Proxy -
type Proxy struct {
	*postgres.Table[*storage.Proxy]
}

// NewProxy -
func NewProxy(db *database.PgGo) *Proxy {
	return &Proxy{
		Table: postgres.NewTable[*storage.Proxy](db),
	}
}

// GetByHash -
func (p *Proxy) GetByHash(ctx context.Context, address, selector []byte) (proxy storage.Proxy, err error) {
	err = p.DB().ModelContext(ctx, &proxy).
		Where("hash = ?", address).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.Where("selector = ?", selector).WhereOr("selector IS NULL")
			return q, nil
		}).
		Limit(1).
		Select(&proxy)
	return
}
