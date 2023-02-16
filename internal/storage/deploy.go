package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IDeploy -
type IDeploy interface {
	storage.Table[*Deploy]
}

// Deploy -
type Deploy struct {
	// nolint
	tableName struct{} `pg:"deploy"`

	ID                  uint64
	BlockID             uint64
	Height              uint64
	Time                int64
	Status              Status `pg:",use_zero"`
	Internal            bool   `pg:",use_zero"`
	Hash                string
	ContractAddressSalt string
	ConstructorCalldata []string
	ClassHash           string
	// TODO: abi entity

	Block Block `pg:"rel:has-one"`
}

// TableName -
func (Deploy) TableName() string {
	return "deploy"
}
