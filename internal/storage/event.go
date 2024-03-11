package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IEvent interface {
	storage.Table[*Event]
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
	bun.BaseModel `bun:"event" comment:"Table with events" partition:"RANGE(time)"`

	ID     uint64    `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height uint64    `comment:"Block height"`
	Time   time.Time `bun:",pk" comment:"Time of block"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`
	FeeID           *uint64 `comment:"Parent fee invocation id"`
	InternalID      *uint64 `comment:"Parent internal transaction id"`

	Order      uint64         `comment:"Order in block"`
	ContractID uint64         `comment:"Contract address id"`
	FromID     uint64         `comment:"From address id"`
	Keys       []string       `bun:",array" comment:"Raw event keys"`
	Data       []string       `bun:",array" comment:"Raw event data"`
	Name       string         `comment:"Event name"`
	ParsedData map[string]any `bun:",nullzero" comment:"Event data parsed according to contract ABI"`

	From     Address `bun:"rel:belongs-to" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	Contract Address `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
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

// Columns -
func (Event) Columns() []string {
	return []string{
		"id", "height", "time", "invoke_id", "declare_id",
		"deploy_id", "deploy_account_id", "l1_handler_id",
		"fee_id", "internal_id", "order", "contract_id",
		"from_id", "keys", "data", "name", "parsed_data",
	}
}

// Flat -
func (e Event) Flat() []any {
	data := []any{
		e.ID,
		e.Height,
		e.Time,
		e.InvokeID,
		e.DeclareID,
		e.DeployID,
		e.DeployAccountID,
		e.L1HandlerID,
		e.FeeID,
		e.InternalID,
		e.Order,
		e.ContractID,
		e.FromID,
		pq.StringArray(e.Keys),
		pq.StringArray(e.Data),
		e.Name,
		nil,
	}

	if e.ParsedData != nil {
		parsed, err := json.MarshalWithOption(e.ParsedData, json.UnorderedMap(), json.DisableNormalizeUTF8())
		if err == nil {
			data[16] = string(parsed)
		}
	}
	return data
}
