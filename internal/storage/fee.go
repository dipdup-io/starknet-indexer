package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IFee -
type IFee interface {
	storage.Table[*Fee]
}

// Fee -
type Fee struct {
	// nolint
	tableName struct{} `pg:"fee,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`
	Status Status    `pg:",use_zero"`

	ContractID     uint64
	CallerID       uint64
	ClassID        uint64
	Selector       []byte
	EntrypointType EntrypointType
	CallType       CallType
	Calldata       []string
	Result         []string

	Entrypoint     string
	ParsedCalldata map[string]any

	InvokeID        *uint64
	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	L1HandlerID     *uint64

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
