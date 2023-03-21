package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Heightable -
type Heightable[T storage.Model] interface {
	ByHeight(ctx context.Context, height, limit, offset uint64) ([]T, error)
}

// Rollback -
type Rollback interface {
	Rollback(ctx context.Context, indexerName string, height uint64) error
}
