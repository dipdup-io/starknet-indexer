package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IFee -
type IFee interface {
	storage.Table[*Fee]
	Copiable[Fee]
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
	// nolint
	tableName struct{} `pg:"fee,partition_by:RANGE(time)" comment:"Table with fee invocations"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height uint64    `pg:",use_zero" comment:"Block height"`
	Time   time.Time `pg:",pk" comment:"Time of block"`

	ContractID uint64 `comment:"Contract address id"`
	CallerID   uint64 `comment:"Caller address id"`
	ClassID    uint64 `comment:"Class id"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`

	EntrypointType EntrypointType `pg:",type:SMALLINT" comment:"Entrypoint type (unknown - 1 | external - 2 | constructor - 3 | l1 handler - 4)"`
	CallType       CallType       `pg:",type:SMALLINT" comment:"Call type (unknwown - 1 | call - 2 | delegate - 3)"`
	Status         Status         `pg:",type:SMALLINT,use_zero" comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`

	Selector       []byte         `comment:"Called selector"`
	Entrypoint     string         `comment:"Entrypoint name"`
	Calldata       []string       `pg:",array" comment:"Raw calldata"`
	Result         []string       `pg:",array" comment:"Raw result"`
	ParsedCalldata map[string]any `comment:"Calldata parsed according to contract ABI"`

	Class     Class      `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Caller    Address    `pg:"rel:has-one" hasura:"table:address,field:caller_id,remote_field:id,type:oto,name:caller"`
	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
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
