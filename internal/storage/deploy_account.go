package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IDeployAccount -
type IDeployAccount interface {
	storage.Table[*DeployAccount]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]DeployAccount, error)
}

// DeployAccount -
type DeployAccount struct {
	// nolint
	tableName struct{} `pg:"deploy_account"`

	ID                  uint64
	Height              uint64 `pg:",use_zero"`
	ClassID             uint64
	ContractID          uint64
	Time                int64
	Status              Status `pg:",use_zero"`
	Hash                []byte
	ContractAddressSalt []byte
	MaxFee              decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce               decimal.Decimal `pg:",type:numeric,use_zero"`
	Signature           []string        `pg:",array"`
	ConstructorCalldata []string        `pg:",array"`
	ParsedCalldata      map[string]any

	Class     Class      `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
}

// TableName -
func (DeployAccount) TableName() string {
	return "deploy_account"
}
