package storage

import (
	"context"
	"io"

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

// Copiable -
type Copiable[M storage.Model] interface {
	InsertByCopy(models []M) (io.Reader, string, error)
}
