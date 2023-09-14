package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IMessage interface {
	storage.Table[*Message]

	Filterable[Message, MessageFilter]
}

// MessageFilter -
type MessageFilter struct {
	ID       IntegerFilter
	Height   IntegerFilter
	Time     TimeFilter
	Contract BytesFilter
	From     BytesFilter
	To       BytesFilter
	Selector EqualityFilter
}

// Message -
type Message struct {
	bun.BaseModel `bun:"message" comment:"Table with messages" partition:"RANGE(time)"`

	ID     uint64    `bun:",pk,autoincrement,nullzero" comment:"Unique internal identity"`
	Height uint64    `comment:"Block height"`
	Time   time.Time `bun:",pk" comment:"Time of block"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`
	FeeID           *uint64 `comment:"Parent fee invocation id"`
	InternalID      *uint64 `comment:"Parent internal transaction id"`

	ContractID uint64          `comment:"Contract address id"`
	Order      uint64          `comment:"Order in block"`
	FromID     uint64          `comment:"From address id"`
	ToID       uint64          `comment:"To address id"`
	Selector   string          `comment:"Called selector"`
	Nonce      decimal.Decimal `bun:",type:numeric" comment:"The transaction nonce"`
	Payload    []string        `bun:",array" comment:"Message payload"`

	From     Address `bun:"rel:belongs-to" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	To       Address `bun:"rel:belongs-to" hasura:"table:address,field:to_id,remote_field:id,type:oto,name:to"`
	Contract Address `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (Message) TableName() string {
	return "message"
}

// GetHeight -
func (msg Message) GetHeight() uint64 {
	return msg.Height
}

// GetId -
func (msg Message) GetId() uint64 {
	return msg.ID
}

// Columns -
func (msg Message) Columns() []string {
	return []string{
		"id", "height", "time", "invoke_id", "declare_id", "deploy_id",
		"deploy_account_id", "l1_handler_id", "fee_id", "internal_id",
		"contract_id", "order", "from_id", "to_id", "selector", "nonce", "payload",
	}
}

// Flat -
func (msg Message) Flat() []any {
	return []any{
		msg.ID,
		msg.Height,
		msg.Time,
		msg.InvokeID,
		msg.DeclareID,
		msg.DeployID,
		msg.DeployAccountID,
		msg.L1HandlerID,
		msg.FeeID,
		msg.InternalID,
		msg.ContractID,
		msg.Order,
		msg.FromID,
		msg.ToID,
		msg.Selector,
		msg.Nonce,
		pq.StringArray(msg.Payload),
	}
}
