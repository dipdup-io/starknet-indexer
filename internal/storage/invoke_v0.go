package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IInvokeV0 -
type IInvokeV0 interface {
	storage.Table[*InvokeV0]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]InvokeV0, error)
}

// InvokeV0 -
type InvokeV0 struct {
	// nolint
	tableName struct{} `pg:"invoke_v0"`

	ID                 uint64
	Height             uint64 `pg:",use_zero"`
	Time               int64
	Status             Status `pg:",use_zero"`
	Hash               []byte
	ContractID         uint64
	EntrypointSelector []byte
	Entrypoint         string
	MaxFee             decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce              decimal.Decimal `pg:",type:numeric,use_zero"`
	Signature          []string        `pg:",array"`
	CallData           []string        `pg:",array"`
	ParsedCalldata     map[string]any

	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
}

// TableName -
func (InvokeV0) TableName() string {
	return "invoke_v0"
}
