package storage

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/types"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IBlock interface {
	storage.Table[*Block]

	ByHeight(ctx context.Context, height uint64) (Block, error)
	Last(ctx context.Context) (Block, error)
	ByStatus(ctx context.Context, status Status, limit, offset uint64, order storage.SortOrder) ([]Block, error)
}

// Block -
type Block struct {
	bun.BaseModel `bun:"block" json:"-" comment:"Block table"`

	ID      uint64    `bun:",pk,autoincrement" json:"id" comment:"Unique internal identity"`
	Height  uint64    `json:"height" comment:"The number (height) of this block"`
	Time    time.Time `json:"time" comment:"The time the sequencer created this block before executing transactions"`
	Version *string   `json:"version" comment:"The version of the Starknet protocol used when creating this block"`

	TxCount            int `json:"tx_count" comment:"Transactions count in block"`
	InvokeCount        int `json:"invoke_count" comment:"Invokes count in block"`
	DeclareCount       int `json:"declare_count" comment:"Declares count in block"`
	DeployCount        int `json:"deploy_count" comment:"Deploys count in block"`
	DeployAccountCount int `json:"deploy_account_count" comment:"Deploy accounts count in block"`
	L1HandlerCount     int `bun:"l1_handler_count" json:"l1_handler_count" comment:"L1 handlers count in block"`
	StorageDiffCount   int `json:"storage_diff_count" comment:"Storage diffs count in block"`

	Status           Status    `json:"status" comment:"Block status"`
	Hash             types.Hex `json:"hash" comment:"Block hash"`
	ParentHash       types.Hex `json:"parent_hash" comment:"The hash of this block's parent"`
	NewRoot          types.Hex `json:"new_root" comment:"The state commitment after this block"`
	SequencerAddress types.Hex `json:"sequencer_address" comment:"The Starknet address of the sequencer who created this block"`

	Invoke        []Invoke        `bun:"rel:has-many" json:"invoke,omitempty"`
	Declare       []Declare       `bun:"rel:has-many" json:"declare,omitempty"`
	Deploy        []Deploy        `bun:"rel:has-many" json:"deploy,omitempty"`
	DeployAccount []DeployAccount `bun:"rel:has-many" json:"deploy_account,omitempty"`
	L1Handler     []L1Handler     `bun:"rel:has-many" json:"l1_handler,omitempty"`
	Fee           []Fee           `bun:"rel:has-many" json:"fee,omitempty"`
	StorageDiffs  []StorageDiff   `bun:"rel:has-many" json:"storage_diffs,omitempty"`
}

// TableName -
func (Block) TableName() string {
	return "block"
}
