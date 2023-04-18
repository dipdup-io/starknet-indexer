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
	tableName struct{} `pg:"block,comment:Block table"`

	ID      uint64    `pg:",comment:Unique internal identity"`
	Height  uint64    `pg:",use_zero,comment:The number (height) of this block"`
	Time    time.Time `pg:",comment:The time the sequencer created this block before executing transactions"`
	Version *string   `pg:",comment:The version of the Starknet protocol used when creating this block"`

	TxCount            int `pg:",use_zero,comment:Transactions count in block"`
	InvokeCount        int `pg:",use_zero,comment:Ivokes count in block"`
	DeclareCount       int `pg:",use_zero,comment:Declares count in block"`
	DeployCount        int `pg:",use_zero,comment:Deploys count in block"`
	DeployAccountCount int `pg:",use_zero,comment:Deploy accounts count in block"`
	L1HandlerCount     int `pg:"l1_handler_count,use_zero,comment:L1 handlers count in block"`
	StorageDiffCount   int `pg:",use_zero,comment:Storage diffs count in block"`

	Status           Status `pg:",comment:Block status"`
	Hash             []byte `pg:",comment:Block hash"`
	ParentHash       []byte `pg:",comment:The hash of this block’s parent"`
	NewRoot          []byte `pg:",comment:The state commitment after this block"`
	SequencerAddress []byte `pg:",comment:The Starknet address of the sequencer who created this block"`

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
