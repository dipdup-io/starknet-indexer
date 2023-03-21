package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IState -
type IState interface {
	storage.Table[*State]

	ByName(ctx context.Context, name string) (State, error)
}

// State -
type State struct {
	// nolint
	tableName struct{} `pg:"state"`

	ID         uint64
	Name       string `pg:",unique:state_name"`
	LastHeight uint64 `pg:",use_zero"`
	LastTime   time.Time

	InvokesCount        uint64 `pg:",use_zero"`
	DeploysCount        uint64 `pg:",use_zero"`
	DeployAccountsCount uint64 `pg:",use_zero"`
	DeclaresCount       uint64 `pg:",use_zero"`
	L1HandlersCount     uint64 `pg:",use_zero"`
	TxCount             uint64 `pg:",use_zero"`

	LastClassID   uint64
	LastAddressID uint64
	LastTxID      uint64
	LastEventID   uint64
}

// TableName -
func (State) TableName() string {
	return "state"
}
