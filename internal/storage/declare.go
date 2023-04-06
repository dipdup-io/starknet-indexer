package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IDeclare -
type IDeclare interface {
	storage.Table[*Declare]

	ByHeight(ctx context.Context, height, limit, offset uint64) ([]Declare, error)
}

// Declare -
type Declare struct {
	// nolint
	tableName struct{} `pg:"declare,partition_by:RANGE(time)"`

	ID         uint64 `pg:"id,type:bigint,pk,notnull"`
	Height     uint64 `pg:",use_zero"`
	ClassID    uint64
	Version    uint64 `pg:",use_zero"`
	Position   int    `pg:",use_zero"`
	SenderID   *uint64
	ContractID *uint64
	Time       time.Time `pg:",pk"`
	Status     Status    `pg:",use_zero"`
	Hash       []byte
	MaxFee     decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce      decimal.Decimal `pg:",type:numeric,use_zero"`

	Signature []string `pg:",array"`

	Class     Class      `pg:"rel:has-one"`
	Sender    Address    `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
}

// TableName -
func (Declare) TableName() string {
	return "declare"
}
