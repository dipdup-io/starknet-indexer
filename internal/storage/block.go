package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IBlock -
type IBlock interface {
	storage.Table[*Block]

	ByHeight(ctx context.Context, height uint64) (Block, error)
	Last(ctx context.Context) (Block, error)
}

// Block -
type Block struct {
	// nolint
	tableName struct{} `pg:"blocks"`

	ID     uint64
	Height uint64
	Time   int64

	TxCount            int `pg:",use_zero"`
	InvokeV0Count      int `pg:",use_zero"`
	InvokeV1Count      int `pg:",use_zero"`
	DeclareCount       int `pg:",use_zero"`
	DeployCount        int `pg:",use_zero"`
	DeployAccountCount int `pg:",use_zero"`
	L1HandlerCount     int `pg:",use_zero"`

	Status           Status
	Hash             string
	ParentHash       string
	NewRoot          string
	SequencerAddress string

	InvokeV0      []InvokeV0      `pg:"rel:has-many"`
	InvokeV1      []InvokeV1      `pg:"rel:has-many"`
	Declare       []Declare       `pg:"rel:has-many"`
	Deploy        []Deploy        `pg:"rel:has-many"`
	DeployAccount []DeployAccount `pg:"rel:has-many"`
	L1Handler     []L1Handler     `pg:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "blocks"
}
