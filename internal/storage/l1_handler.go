package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IL1Handler -
type IL1Handler interface {
	storage.Table[*L1Handler]
}

// L1Handler -
type L1Handler struct {
	// nolint
	tableName struct{} `pg:"l1_handler"`

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
func (L1Handler) TableName() string {
	return "l1_handler"
}
