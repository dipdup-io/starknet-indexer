package storage

import (
	"context"
	"io"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Rollback -
type Rollback interface {
	Rollback(ctx context.Context, indexerName string, height uint64) error
}

// Copiable -
type Copiable[M storage.Model] interface {
	InsertByCopy(models []M) (io.Reader, string, error)
}

// Heightable -
type Heightable interface {
	storage.Model

	GetHeight() uint64
}

// Filterable -
type Filterable[M storage.Model, F any] interface {
	Filter(ctx context.Context, flt F, opts ...FilterOption) ([]M, error)
}

// FilterOptions -
type FilterOptions struct {
	Limit  int
	Offset int

	SortField string
	SortOrder storage.SortOrder
}

// FilterOption -
type FilterOption func(opt *FilterOptions)

// WithLimitFilter -
func WithLimitFilter(limit int) FilterOption {
	return func(opt *FilterOptions) {
		if limit > 0 {
			opt.Limit = limit
		}
	}
}

// WithOffsetFilter -
func WithOffsetFilter(offset int) FilterOption {
	return func(opt *FilterOptions) {
		if offset > 0 {
			opt.Offset = offset
		}
	}
}

// WithSortFilter -
func WithSortFilter(field string, order storage.SortOrder) FilterOption {
	return func(opt *FilterOptions) {
		opt.SortField = field
		opt.SortOrder = order
	}
}

// WithAscSortByIdFilter -
func WithAscSortByIdFilter() FilterOption {
	return func(opt *FilterOptions) {
		opt.SortField = "id"
		opt.SortOrder = storage.SortOrderAsc
	}
}

// Models - list all models
var Models = []storage.Model{
	&State{},
	&Address{},
	&Class{},
	&StorageDiff{},
	&Block{},
	&Invoke{},
	&Declare{},
	&Deploy{},
	&DeployAccount{},
	&L1Handler{},
	&Internal{},
	&Event{},
	&Message{},
	&Transfer{},
	&Fee{},
	&ERC20{},
	&ERC721{},
	&ERC1155{},
	&TokenBalance{},
	&Proxy{},
}
