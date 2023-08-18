package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Rollback -
type Rollback interface {
	Rollback(ctx context.Context, indexerName string, height uint64) error
}

// CopiableModel -
type CopiableModel interface {
	storage.Model
	storage.Copiable
}

// Heightable -
type Heightable interface {
	storage.Model

	GetHeight() uint64
	GetId() uint64
}

// Filterable -
type Filterable[M storage.Model, F any] interface {
	Filter(ctx context.Context, flt []F, opts ...FilterOption) ([]M, error)
}

// FilterOptions -
type FilterOptions struct {
	Limit  int
	Offset int

	SortField string
	SortOrder storage.SortOrder

	MaxHeight        uint64
	HeightColumnName string
	Cursor           uint64
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

// WithDescSortByIdFilter -
func WithDescSortByIdFilter() FilterOption {
	return func(opt *FilterOptions) {
		opt.SortField = "id"
		opt.SortOrder = storage.SortOrderDesc
	}
}

// WithMaxHeight -
func WithMaxHeight(height uint64, columnName string) FilterOption {
	return func(opt *FilterOptions) {
		opt.MaxHeight = height
		opt.HeightColumnName = columnName
	}
}

// WithCursor -
func WithCursor(id uint64) FilterOption {
	return func(opt *FilterOptions) {
		opt.Cursor = id
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
	&Token{},
	&TokenBalance{},
	&ProxyUpgrade{},
	&Proxy{},
}

// ModelsAny - list all models with any type
var ModelsAny = []any{
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
	&Token{},
	&TokenBalance{},
	&ProxyUpgrade{},
	&Proxy{},
}
