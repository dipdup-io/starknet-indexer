package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IFee interface {
	storage.Table[*Fee]
	Filterable[Fee, FeeFilter]
}

// FeeFilter -
type FeeFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Contract       BytesFilter
	Caller         BytesFilter
	Class          BytesFilter
	Selector       EqualityFilter
	Entrypoint     StringFilter
	EntrypointType EnumFilter
	CallType       EnumFilter
	ParsedCalldata map[string]string
}

// Fee -
type Fee struct {
	bun.BaseModel `bun:"fee" comment:"Table with fee invocations" partition:"RANGE(time)"`

	ID     uint64    `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height uint64    `comment:"Block height"`
	Time   time.Time `bun:",pk" comment:"Time of block"`

	ContractID uint64 `comment:"Contract address id"`
	CallerID   uint64 `comment:"Caller address id"`
	ClassID    uint64 `comment:"Class id"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`

	EntrypointType EntrypointType `bun:",type:SMALLINT" comment:"Entrypoint type (unknown - 1 | external - 2 | constructor - 3 | l1 handler - 4)"`
	CallType       CallType       `bun:",type:SMALLINT" comment:"Call type (unknwown - 1 | call - 2 | delegate - 3)"`
	Status         Status         `bun:",type:SMALLINT" comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`

	Selector       []byte         `comment:"Called selector"`
	Entrypoint     string         `comment:"Entrypoint name"`
	Calldata       []string       `bun:",array" comment:"Raw calldata"`
	Result         []string       `bun:",array" comment:"Raw result"`
	ParsedCalldata map[string]any `bun:",nullzero" comment:"Calldata parsed according to contract ABI"`

	Class     Class      `bun:"rel:belongs-to" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Caller    Address    `bun:"rel:belongs-to" hasura:"table:address,field:caller_id,remote_field:id,type:oto,name:caller"`
	Contract  Address    `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many"`
	Messages  []Message  `bun:"rel:has-many"`
	Events    []Event    `bun:"rel:has-many"`
	Transfers []Transfer `bun:"rel:has-many"`
}

// TableName -
func (Fee) TableName() string {
	return "fee"
}

// GetHeight -
func (f Fee) GetHeight() uint64 {
	return f.Height
}

// GetId -
func (f Fee) GetId() uint64 {
	return f.ID
}

// Columns -
func (Fee) Columns() []string {
	return []string{
		"id", "height", "time", "contract_id", "caller_id",
		"class_id", "invoke_id", "declare_id",
		"deploy_id", "deploy_account_id", "l1_handler_id",
		"entrypoint_type", "call_type", "status", "selector",
		"entrypoint", "calldata", "result", "parsed_calldata",
	}
}

// Flat -
func (f Fee) Flat() []any {
	data := []any{
		f.ID,
		f.Height,
		f.Time,
		f.ContractID,
		f.CallerID,
		f.ClassID,
		f.InvokeID,
		f.DeclareID,
		f.DeployID,
		f.DeployAccountID,
		f.L1HandlerID,
		f.EntrypointType,
		f.CallType,
		f.Status,
		f.Selector,
		f.Entrypoint,
		pq.StringArray(f.Calldata),
		pq.StringArray(f.Result),
		nil,
	}

	if f.ParsedCalldata != nil {
		parsed, err := json.MarshalWithOption(f.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
		if err == nil {
			data[18] = string(parsed)
		}
	}
	return data
}
