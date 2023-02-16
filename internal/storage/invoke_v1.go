package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IInvokeV1 -
type IInvokeV1 interface {
	storage.Table[*InvokeV1]
}

// InvokeV1 -
type InvokeV1 struct {
	// nolint
	tableName struct{} `pg:"invoke_v1"`

	ID            uint64
	BlockID       uint64
	Height        uint64
	Time          int64
	Status        Status `pg:",use_zero"`
	Internal      bool   `pg:",use_zero"`
	Hash          string
	SenderAddress string
	MaxFee        decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce         decimal.Decimal `pg:",use_zero"`
	Signature     []string
	CallData      []string

	Block Block `pg:"rel:has-one"`
}

// TableName -
func (InvokeV1) TableName() string {
	return "invoke_v1"
}
