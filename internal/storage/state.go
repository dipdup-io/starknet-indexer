package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

// IState -
type IState interface {
	storage.Table[*State]

	ByName(ctx context.Context, name string) (State, error)
}

// State -
type State struct {
	bun.BaseModel `bun:"state"`

	ID         uint64    `bun:",pk,autoincrement" comment:"Unique internal identity"`
	Name       string    `bun:",unique:state_name"`
	LastHeight uint64    `comment:"Last block height"`
	LastTime   time.Time `comment:"Time of last block"`

	InvokesCount        uint64 `comment:"Total invokes count"`
	DeploysCount        uint64 `comment:"Total deploys count"`
	DeployAccountsCount uint64 `comment:"Total deploy accounts count"`
	DeclaresCount       uint64 `comment:"Total declares count"`
	L1HandlersCount     uint64 `comment:"Total l1 handlers count"`
	TxCount             uint64 `comment:"Total transactions count"`

	LastClassID   uint64 `comment:"Last internal class id"`
	LastAddressID uint64 `comment:"Last internal address id"`
	LastTxID      uint64 `comment:"Last internal transaction id"`
	LastEventID   uint64 `comment:"Last internal event id"`
}

// TableName -
func (State) TableName() string {
	return "state"
}
