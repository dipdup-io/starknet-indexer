package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IClass -
type IClass interface {
	storage.Table[*Class]

	GetByHash(ctx context.Context, hash []byte) (Class, error)
}

// Class -
type Class struct {
	// nolint
	tableName struct{} `pg:"class"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Type   ClassType `pg:",use_zero"`
	Hash   []byte    `pg:",unique:class_hash"`
	Abi    Bytes     `pg:",type:bytea"`
	Height uint64    `pg:",use_zero"`
}

// TableName -
func (Class) TableName() string {
	return "class"
}
