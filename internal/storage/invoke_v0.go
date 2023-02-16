package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IInvokeV0 -
type IInvokeV0 interface {
	storage.Table[*InvokeV0]
}

// InvokeV0 -
type InvokeV0 struct {
	// nolint
	tableName struct{} `pg:"invoke_v0"`

	ID                 uint64
	BlockID            uint64
	Height             uint64
	Time               int64
	Status             Status `pg:",use_zero"`
	Internal           bool   `pg:",use_zero"`
	Hash               string
	ContractAddress    string
	EntrypointSelector string
	MaxFee             decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce              decimal.Decimal `pg:",type:numeric,use_zero"`
	Signature          []string
	CallData           []string

	Block Block `pg:"rel:has-one"`
}

// TableName -
func (InvokeV0) TableName() string {
	return "invoke_v0"
}
