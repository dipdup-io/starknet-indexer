package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Proxy -
type Proxy struct {
	*postgres.Table[*storage.Proxy]
}

// NewProxy -
func NewProxy(db *database.Bun) *Proxy {
	return &Proxy{
		Table: postgres.NewTable[*storage.Proxy](db),
	}
}

// GetByHash -
func (p *Proxy) GetByHash(ctx context.Context, address, selector []byte) (proxy storage.Proxy, err error) {
	err = p.DB().NewSelect().Model(&proxy).
		Where("hash = ?", address).
		WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			q = q.Where("selector IS NULL")
			if len(selector) > 0 {
				q = q.WhereOr("selector = ?", selector)
			}
			return q
		}).
		Limit(1).
		Scan(ctx)
	return
}
