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

	ID         uint64    `comment:"Unique internal identity"`
	Name       string    `pg:",unique:state_name"`
	LastHeight uint64    `pg:",use_zero" comment:"Last block height"`
	LastTime   time.Time `comment:"Time of last block"`

	InvokesCount        uint64 `pg:",use_zero" comment:"Total invokes count"`
	DeploysCount        uint64 `pg:",use_zero" comment:"Total deploys count"`
	DeployAccountsCount uint64 `pg:",use_zero" comment:"Total deploy accounts count"`
	DeclaresCount       uint64 `pg:",use_zero" comment:"Total declares count"`
	L1HandlersCount     uint64 `pg:",use_zero" comment:"Total l1 handlers count"`
	TxCount             uint64 `pg:",use_zero" comment:"Total transactions count"`

	LastClassID   uint64 `comment:"Last internal class id"`
	LastAddressID uint64 `comment:"Last internal address id"`
	LastTxID      uint64 `comment:"Last internal transaction id"`
	LastEventID   uint64 `comment:"Last internal event id"`
}

// TableName -
func (State) TableName() string {
	return "state"
}
