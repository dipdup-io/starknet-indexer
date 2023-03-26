package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IBlock -
type IBlock interface {
	storage.Table[*Block]

	ByHeight(ctx context.Context, height uint64) (Block, error)
	Last(ctx context.Context) (Block, error)
	ByStatus(ctx context.Context, status Status, limit, offset uint64, order storage.SortOrder) ([]Block, error)
}

// Block -
type Block struct {
	// nolint
	tableName struct{} `pg:"block"`

	ID      uint64
	Height  uint64 `pg:",use_zero"`
	Time    time.Time
	Version *string

	TxCount            int `pg:",use_zero"`
	InvokeCount        int `pg:",use_zero"`
	DeclareCount       int `pg:",use_zero"`
	DeployCount        int `pg:",use_zero"`
	DeployAccountCount int `pg:",use_zero"`
	L1HandlerCount     int `pg:",use_zero"`
	StorageDiffCount   int `pg:",use_zero"`

	Status           Status
	Hash             []byte
	ParentHash       []byte
	NewRoot          []byte
	SequencerAddress []byte

	Invoke        []Invoke        `pg:"rel:has-many"`
	Declare       []Declare       `pg:"rel:has-many"`
	Deploy        []Deploy        `pg:"rel:has-many"`
	DeployAccount []DeployAccount `pg:"rel:has-many"`
	L1Handler     []L1Handler     `pg:"rel:has-many"`
	Fee           []Fee           `pg:"rel:has-many"`
	StorageDiffs  []StorageDiff   `pg:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
