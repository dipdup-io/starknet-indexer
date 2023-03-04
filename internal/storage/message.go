package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IMessage -
type IMessage interface {
	storage.Table[*Message]
}

// Message -
type Message struct {
	// nolint
	tableName struct{} `pg:"message"`

	ID     uint64
	Height uint64 `pg:",use_zero"`
	Time   int64

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeV0ID      *uint64 `pg:"invoke_v0_id"`
	InvokeV1ID      *uint64 `pg:"invoke_v1_id"`
	L1HandlerID     *uint64
	InternalID      *uint64

	Order    uint64
	FromID   uint64
	ToID     uint64
	Selector string
	Nonce    decimal.Decimal `pg:",type:numeric,use_zero"`
	Payload  []string        `pg:",array"`

	From Address `pg:"rel:has-one"`
	To   Address `pg:"rel:has-one"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}
