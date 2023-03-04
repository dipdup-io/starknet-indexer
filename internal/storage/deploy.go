package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IDeploy -
type IDeploy interface {
	storage.Table[*Deploy]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]Deploy, error)
}

// Deploy -
type Deploy struct {
	// nolint
	tableName struct{} `pg:"deploy"`

	ID                  uint64
	Height              uint64 `pg:",use_zero"`
	ClassID             uint64
	ContractID          uint64
	Time                int64
	Status              Status `pg:",use_zero"`
	Hash                []byte
	ContractAddressSalt []byte
	ConstructorCalldata []string `pg:",array"`
	ParsedCalldata      map[string]any

	Class     Class      `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
}

// TableName -
func (Deploy) TableName() string {
	return "deploy"
}
