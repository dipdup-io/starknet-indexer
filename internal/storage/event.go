package storage

import (
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IEvent -
type IEvent interface {
	storage.Table[*Event]
}

// Event -
type Event struct {
	// nolint
	tableName struct{} `pg:"event"`

	ID     uint64
	Height uint64 `pg:",use_zero"`
	Time   int64

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeV0ID      *uint64 `pg:"invoke_v0_id"`
	InvokeV1ID      *uint64 `pg:"invoke_v1_id"`
	L1HandlerID     *uint64
	InternalID      *uint64

	Order      uint64
	FromID     uint64
	Keys       []string `pg:",array"`
	Data       []string `pg:",array"`
	Name       string
	ParsedData map[string]any

	From Address `pg:"rel:has-one"`
}

// TableName -
func (Event) TableName() string {
	return "event"
}
