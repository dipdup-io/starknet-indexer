package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IL1Handler -
type IL1Handler interface {
	storage.Table[*L1Handler]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]L1Handler, error)
}

// L1Handler -
type L1Handler struct {
	// nolint
	tableName struct{} `pg:"l1_handler"`

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
func (L1Handler) TableName() string {
	return "l1_handler"
}
