package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// ClassReplace -
type ClassReplace struct {
	*postgres.Table[*storage.ClassReplace]
}

// NewClassReplace -
func NewClassReplace(db *database.Bun) *ClassReplace {
	return &ClassReplace{
		Table: postgres.NewTable[*storage.ClassReplace](db),
	}
}

func (cr *ClassReplace) ByHeight(ctx context.Context, height uint64) (replaces []storage.ClassReplace, err error) {
	err = cr.DB().NewSelect().Model(&replaces).Where("height = ?", height).Scan(ctx)
	return
}
