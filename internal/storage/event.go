package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IEvent -
type IEvent interface {
	storage.Table[*Event]
}

// Event -
type Event struct {
	// nolint
	tableName struct{} `pg:"event,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeID        *uint64
	L1HandlerID     *uint64
	FeeID           *uint64
	InternalID      *uint64

	Order      uint64
	ContractID uint64
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
