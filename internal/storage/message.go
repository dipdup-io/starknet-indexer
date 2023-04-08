package storage

import (
	"time"

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
	tableName struct{} `pg:"message,partition_by:RANGE(time)"`

	ID     uint64    `pg:",pk"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeID        *uint64
	L1HandlerID     *uint64
	FeeID           *uint64
	InternalID      *uint64

	ContractID uint64
	Order      uint64
	FromID     uint64
	ToID       uint64
	Selector   string
	Nonce      decimal.Decimal `pg:",type:numeric,use_zero"`
	Payload    []string        `pg:",array"`

	From     Address `pg:"rel:has-one"`
	To       Address `pg:"rel:has-one"`
	Contract Address `pg:"rel:has-one"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}
