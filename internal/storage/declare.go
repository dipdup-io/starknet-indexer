package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IDeclare -
type IDeclare interface {
	storage.Table[*Declare]
}

// Declare -
type Declare struct {
	// nolint
	tableName struct{} `pg:"declare"`

	ID            uint64
	BlockID       uint64
	Height        uint64
	Time          int64
	Status        Status `pg:",use_zero"`
	Internal      bool   `pg:",use_zero"`
	Hash          string
	SenderAddress string
	ClassHash     string
	MaxFee        decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce         decimal.Decimal `pg:",type:numeric,use_zero"`
	// TODO: abi entity
	Signature []string

	Block Block `pg:"rel:has-one"`
}

// TableName -
func (Declare) TableName() string {
	return "declare"
}
