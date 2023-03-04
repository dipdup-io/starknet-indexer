package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IInvokeV1 -
type IInvokeV1 interface {
	storage.Table[*InvokeV1]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]InvokeV1, error)
}

// InvokeV1 -
type InvokeV1 struct {
	// nolint
	tableName struct{} `pg:"invoke_v1"`

	ID             uint64
	Height         uint64 `pg:",use_zero"`
	Time           int64
	Status         Status `pg:",use_zero"`
	Hash           []byte
	SenderID       uint64
	ContractID     uint64
	MaxFee         decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce          decimal.Decimal `pg:",type:numeric,use_zero"`
	Signature      []string        `pg:",array"`
	CallData       []string        `pg:",array"`
	ParsedCalldata map[string]any

	Sender    Address    `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
}

// TableName -
func (InvokeV1) TableName() string {
	return "invoke_v1"
}
