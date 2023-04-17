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
	tableName struct{} `pg:"fee,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`

	ContractID uint64
	CallerID   uint64
	ClassID    uint64

	InvokeID        *uint64
	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	L1HandlerID     *uint64

	EntrypointType EntrypointType `pg:"type:SMALLINT"`
	CallType       CallType       `pg:"type:SMALLINT"`
	Status         Status         `pg:"type:SMALLINT,use_zero"`

	Selector       []byte
	Entrypoint     string
	Calldata       []string `pg:",array"`
	Result         []string `pg:",array"`
	ParsedCalldata map[string]any

	Class     Class      `pg:"rel:has-one"`
	Caller    Address    `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
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
