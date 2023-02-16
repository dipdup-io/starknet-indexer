package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IDeployAccount -
type IDeployAccount interface {
	storage.Table[*DeployAccount]
}

// DeployAccount -
type DeployAccount struct {
	// nolint
	tableName struct{} `pg:"deploy_account"`

	ID                  uint64
	BlockID             uint64
	Height              uint64
	Time                int64
	Status              Status `pg:",use_zero"`
	Internal            bool   `pg:",use_zero"`
	Hash                string
	ContractAddressSalt string
	ClassHash           string
	MaxFee              decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce               decimal.Decimal `pg:",type:numeric,use_zero"`
	ConstructorCalldata []string
	Signature           []string

	Block Block `pg:"rel:has-one"`
}

// TableName -
func (DeployAccount) TableName() string {
	return "deploy_account"
}
