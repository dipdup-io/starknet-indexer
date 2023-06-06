package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IEvent -
type IEvent interface {
	storage.Table[*Event]

	Copiable[Event]
	Filterable[Event, EventFilter]
}

// EventFilter -
type EventFilter struct {
	ID         IntegerFilter
	Height     IntegerFilter
	Time       TimeFilter
	Contract   IdFilter
	From       IdFilter
	Name       StringFilter
	ParsedData map[string]string
}

// Event -
type Event struct {
	// nolint
	tableName struct{} `pg:"event,partition_by:RANGE(time),comment:Table with events"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height uint64    `pg:",use_zero,comment:Block height"`
	Time   time.Time `pg:",pk,comment:Time of block"`

	InvokeID        *uint64 `pg:",comment:Parent invoke id"`
	DeclareID       *uint64 `pg:",comment:Parent declare id"`
	DeployID        *uint64 `pg:",comment:Parent deploy id"`
	DeployAccountID *uint64 `pg:",comment:Parent deploy account id"`
	L1HandlerID     *uint64 `pg:",comment:Parent l1 handler id"`
	FeeID           *uint64 `pg:",comment:Parent fee invocation id"`
	InternalID      *uint64 `pg:",comment:Parent internal transaction id"`

	Order      uint64         `pg:",comment:Order in block"`
	ContractID uint64         `pg:",comment:Contract address id"`
	FromID     uint64         `pg:",comment:From address id"`
	Keys       []string       `pg:",array,comment:Raw event keys"`
	Data       []string       `pg:",array,comment:Raw event data"`
	Name       string         `pg:",comment:Event name"`
	ParsedData map[string]any `pg:",comment:Event data parsed according to contract ABI"`

	From     Address `pg:"rel:has-one" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}

// GetHeight -
func (e Event) GetHeight() uint64 {
	return e.Height
}

// GetId -
func (e Event) GetId() uint64 {
	return e.ID
}
