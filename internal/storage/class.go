package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IClass interface {
	storage.Table[*Class]

	GetByHash(ctx context.Context, hash []byte) (Class, error)
	GetUnresolved(ctx context.Context) ([]Class, error)
}

// Class -
type Class struct {
	bun.BaseModel `bun:"class" comment:"Classes table"`

	ID     uint64    `bun:"id,type:bigint,pk,notnull,nullzero" comment:"Unique internal identity"`
	Type   ClassType `comment:"Class type. It’s a binary mask."`
	Hash   []byte    `bun:",unique:class_hash" comment:"Class hash"`
	Abi    Bytes     `bun:",type:bytea" comment:"Class abi in a raw"`
	Height uint64    `comment:"Block height of the first class occurance"`
	Cairo  int       `bun:",default:0,type:SMALLINT" comment:"Cairo version of class"`
}

// TableName -
func (Class) TableName() string {
	return "class"
}
