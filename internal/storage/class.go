package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IClass -
type IClass interface {
	storage.Table[*Class]

	GetByHash(ctx context.Context, hash []byte) (Class, error)
	GetUnresolved(ctx context.Context) ([]Class, error)
}

// Class -
type Class struct {
	// nolint
	tableName struct{} `pg:"class,comment:Classes table"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Type   ClassType `pg:",use_zero,comment:Class type. It’s a binary mask."`
	Hash   []byte    `pg:",unique:class_hash,comment:Class hash"`
	Abi    Bytes     `pg:",type:bytea,comment:Class abi in a raw"`
	Height uint64    `pg:",use_zero,comment:Block height of the first class occurance"`
	Cairo  int       `pg:",default:0,type:SMALLINT,comment:Cairo version of class"`
}

// TableName -
func (Class) TableName() string {
	return "class"
}
