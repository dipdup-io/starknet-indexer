package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IMessage -
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
	// nolint
	tableName struct{} `pg:"message,partition_by:RANGE(time),comment:Table with messages"`

	ID     uint64    `pg:",pk,comment:Unique internal identity"`
	Height uint64    `pg:",use_zero,comment:Block height"`
	Time   time.Time `pg:",pk,comment:Time of block"`

	InvokeID        *uint64 `pg:",comment:Parent invoke id"`
	DeclareID       *uint64 `pg:",comment:Parent declare id"`
	DeployID        *uint64 `pg:",comment:Parent deploy id"`
	DeployAccountID *uint64 `pg:",comment:Parent deploy account id"`
	L1HandlerID     *uint64 `pg:",comment:Parent l1 handler id"`
	FeeID           *uint64 `pg:",comment:Parent fee invocation id"`
	InternalID      *uint64 `pg:",comment:Parent internal transaction id"`

	ContractID uint64          `pg:",comment:Contract address id"`
	Order      uint64          `pg:",comment:Order in block"`
	FromID     uint64          `pg:",comment:From address id"`
	ToID       uint64          `pg:",comment:To address id"`
	Selector   string          `pg:",comment:Called selector"`
	Nonce      decimal.Decimal `pg:",type:numeric,use_zero,comment:The transaction nonce"`
	Payload    []string        `pg:",array,comment:Message payload"`

	From     Address `pg:"rel:has-one" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	To       Address `pg:"rel:has-one" hasura:"table:address,field:to_id,remote_field:id,type:oto,name:to"`
	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
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
