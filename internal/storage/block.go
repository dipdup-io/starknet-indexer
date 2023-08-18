package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
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
	bun.BaseModel `bun:"block" comment:"Block table"`

	ID      uint64    `bun:",pk,autoincrement" comment:"Unique internal identity"`
	Height  uint64    `comment:"The number (height) of this block"`
	Time    time.Time `comment:"The time the sequencer created this block before executing transactions"`
	Version *string   `comment:"The version of the Starknet protocol used when creating this block"`

	TxCount            int `comment:"Transactions count in block"`
	InvokeCount        int `comment:"Ivokes count in block"`
	DeclareCount       int `comment:"Declares count in block"`
	DeployCount        int `comment:"Deploys count in block"`
	DeployAccountCount int `comment:"Deploy accounts count in block"`
	L1HandlerCount     int `bun:"l1_handler_count" comment:"L1 handlers count in block"`
	StorageDiffCount   int `comment:"Storage diffs count in block"`

	Status           Status `comment:"Block status"`
	Hash             []byte `comment:"Block hash"`
	ParentHash       []byte `comment:"The hash of this block’s parent"`
	NewRoot          []byte `comment:"The state commitment after this block"`
	SequencerAddress []byte `comment:"The Starknet address of the sequencer who created this block"`

	Invoke        []Invoke        `bun:"rel:has-many"`
	Declare       []Declare       `bun:"rel:has-many"`
	Deploy        []Deploy        `bun:"rel:has-many"`
	DeployAccount []DeployAccount `bun:"rel:has-many"`
	L1Handler     []L1Handler     `bun:"rel:has-many"`
	Fee           []Fee           `bun:"rel:has-many"`
	StorageDiffs  []StorageDiff   `bun:"rel:has-many"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
