package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IInternal -
type IInternal interface {
	storage.Table[*Internal]
}

// Internal -
type Internal struct {
	// nolint
	tableName struct{} `pg:"internal_tx"`

	ID     uint64
	Height uint64 `pg:",use_zero"`
	Time   int64
	Status Status `pg:",use_zero"`
	Hash   []byte

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeV0ID      *uint64 `pg:"invoke_v0_id"`
	InvokeV1ID      *uint64 `pg:"invoke_v1_id"`
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
}

// TableName -
func (Internal) TableName() string {
	return "internal_tx"
}
