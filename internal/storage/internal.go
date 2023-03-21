package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IInternal -
type IInternal interface {
	storage.Table[*Internal]
}

// Internal -
type Internal struct {
	// nolint
	tableName struct{} `pg:"internal_tx,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`
	Status Status    `pg:",use_zero"`
	Hash   []byte

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeID        *uint64
	L1HandlerID     *uint64
	InternalID      *uint64

	ClassID        uint64
	CallerID       uint64
	ContractID     uint64
	CallType       CallType
	EntrypointType EntrypointType
	Selector       []byte
	Entrypoint     string
	Result         []string `pg:",array"`
	Calldata       []string `pg:",array"`
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
func (Internal) TableName() string {
	return "internal_tx"
}
