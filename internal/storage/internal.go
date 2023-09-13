package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IInternal interface {
	storage.Table[*Internal]
	Filterable[Internal, InternalFilter]
}

// InternalFilter -
type InternalFilter struct {
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

// Internal -
type Internal struct {
	bun.BaseModel `bun:"internal_tx" comment:"Table with internal transactions" partition:"RANGE(time)"`

	ID     uint64    `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height uint64    `comment:"Block height"`
	Time   time.Time `bun:",pk" comment:"Time of block"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`
	InternalID      *uint64 `comment:"Parent internal transaction id"`
	ClassID         uint64  `comment:"Class id"`
	CallerID        uint64  `comment:"Caller address id"`
	ContractID      uint64  `comment:"Contract address id"`

	Status         Status         `bun:",type:SMALLINT" comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	EntrypointType EntrypointType `bun:",type:SMALLINT" comment:"Entrypoint type (unknown - 1 | external - 2 | constructor - 3 | l1 handler - 4)"`
	CallType       CallType       `bun:",type:SMALLINT" comment:"Call type (unknwown - 1 | call - 2 | delegate - 3)"`

	Hash           []byte         `comment:"Transaction hash"`
	Selector       []byte         `comment:"Called selector"`
	Entrypoint     string         `comment:"Entrypoint name"`
	Result         []string       `bun:",array" comment:"Raw result"`
	Calldata       []string       `bun:",array" comment:"Raw calldata"`
	ParsedCalldata map[string]any `comment:"Calldata parsed according to contract ABI"`
	ParsedResult   map[string]any `comment:"Result parsed according to contract ABI"`

	Class     Class      `bun:"rel:belongs-to" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Caller    Address    `bun:"rel:belongs-to" hasura:"table:address,field:caller_id,remote_field:id,type:oto,name:caller"`
	Contract  Address    `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many"`
	Messages  []Message  `bun:"rel:has-many"`
	Events    []Event    `bun:"rel:has-many"`
	Transfers []Transfer `bun:"rel:has-many"`
}

// TableName -
func (Internal) TableName() string {
	return "internal_tx"
}

// GetHeight -
func (i Internal) GetHeight() uint64 {
	return i.Height
}

// GetId -
func (i Internal) GetId() uint64 {
	return i.ID
}

// Columns -
func (Internal) Columns() []string {
	return []string{
		"id", "height", "time", "contract_id", "caller_id",
		"class_id", "invoke_id", "declare_id", "internal_id",
		"deploy_id", "deploy_account_id", "l1_handler_id",
		"entrypoint_type", "call_type", "status", "selector",
		"hash", "entrypoint", "calldata", "result",
		"parsed_calldata", "parsed_result",
	}
}

// Flat -
func (i Internal) Flat() []any {
	data := []any{
		i.ID,
		i.Height,
		i.Time,
		i.ContractID,
		i.CallerID,
		i.ClassID,
		i.InvokeID,
		i.DeclareID,
		i.InternalID,
		i.DeployID,
		i.DeployAccountID,
		i.L1HandlerID,
		i.EntrypointType,
		i.CallType,
		i.Status,
		i.Selector,
		i.Hash,
		i.Entrypoint,
		pq.StringArray(i.Calldata),
		pq.StringArray(i.Result),
	}

	parsed, err := json.MarshalWithOption(i.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
	if err != nil {
		data = append(data, nil)
	} else {
		data = append(data, string(parsed))
	}

	result, err := json.MarshalWithOption(i.ParsedResult, json.UnorderedMap(), json.DisableNormalizeUTF8())
	if err != nil {
		data = append(data, nil)
	} else {
		data = append(data, string(result))
	}
	return data
}
