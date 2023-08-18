package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// ProxyUpgrade -
type ProxyUpgrade struct {
	*postgres.Table[*storage.ProxyUpgrade]
}

// NewProxyUpgrade -
func NewProxyUpgrade(db *database.Bun) *ProxyUpgrade {
	return &ProxyUpgrade{
		Table: postgres.NewTable[*storage.ProxyUpgrade](db),
	}
}

// LastBefore -
func (pu *ProxyUpgrade) LastBefore(ctx context.Context, hash, selector []byte, height uint64) (upg storage.ProxyUpgrade, err error) {
	query := pu.DB().NewSelect().Model(&upg).
		Where("height < ?", height).
		Where("hash = ?", hash).
		Limit(1).
		Order("id desc")

	if len(selector) == 0 {
		query = query.Where("selector IS NULL")
	} else {
		query = query.Where("selector = ?", selector)
	}

	err = query.Scan(ctx)
	return
}

// ListWithHeight -
func (pu *ProxyUpgrade) ListWithHeight(ctx context.Context, height uint64, limit, offset int) (upgrades []storage.ProxyUpgrade, err error) {
	if limit == 0 {
		limit = 10
	}

	err = pu.DB().NewSelect().Model(&upgrades).
		Where("height > ?", height).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Scan(ctx)

	return upgrades, err
}
