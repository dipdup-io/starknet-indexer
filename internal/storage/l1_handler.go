package storage

import (
	"context"
	"time"

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
	tableName struct{} `pg:"l1_handler,partition_by:RANGE(time)"`

	ID                 uint64    `pg:"id,type:bigint,pk,notnull"`
	Height             uint64    `pg:",use_zero"`
	Time               time.Time `pg:",pk"`
	Status             Status    `pg:",use_zero"`
	Hash               []byte
	ContractID         uint64
	Position           int `pg:",use_zero"`
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
	Transfers []Transfer `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
}

// TableName -
func (L1Handler) TableName() string {
	return "l1_handler"
}
