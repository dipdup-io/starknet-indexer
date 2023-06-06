package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IInternal -
type IInternal interface {
	storage.Table[*Internal]

	Copiable[Internal]
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
	// nolint
	tableName struct{} `pg:"internal_tx,partition_by:RANGE(time),comment:Table with internal transactions"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height uint64    `pg:",use_zero,comment:Block height"`
	Time   time.Time `pg:",pk,comment:Time of block"`

	InvokeID        *uint64 `pg:",comment:Parent invoke id"`
	DeclareID       *uint64 `pg:",comment:Parent declare id"`
	DeployID        *uint64 `pg:",comment:Parent deploy id"`
	DeployAccountID *uint64 `pg:",comment:Parent deploy account id"`
	L1HandlerID     *uint64 `pg:",comment:Parent l1 handler id"`
	InternalID      *uint64 `pg:",comment:Parent internal transaction id"`
	ClassID         uint64  `pg:",comment:Class id"`
	CallerID        uint64  `pg:",comment:Caller address id"`
	ContractID      uint64  `pg:",comment:Contract address id"`

	Status         Status         `pg:"type:SMALLINT,use_zero,comment:Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	EntrypointType EntrypointType `pg:"type:SMALLINT,comment:Entrypoint type (unknown - 1 | external - 2 | constructor - 3 | l1 handler - 4)"`
	CallType       CallType       `pg:"type:SMALLINT,comment:Call type (unknwown - 1 | call - 2 | delegate - 3)"`

	Hash           []byte         `pg:",comment:Transaction hash"`
	Selector       []byte         `pg:",comment:Called selector"`
	Entrypoint     string         `pg:",comment:Entrypoint name"`
	Result         []string       `pg:",array,comment:Raw result"`
	Calldata       []string       `pg:",array,comment:Raw calldata"`
	ParsedCalldata map[string]any `pg:",comment:Calldata parsed according to contract ABI"`
	ParsedResult   map[string]any `pg:",comment:Result parsed according to contract ABI"`

	Class     Class      `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Caller    Address    `pg:"rel:has-one" hasura:"table:address,field:caller_id,remote_field:id,type:oto,name:caller"`
	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Token     *Token     `pg:"-"`
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
