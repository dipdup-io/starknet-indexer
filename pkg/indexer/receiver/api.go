package receiver

import (
	"context"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

type API interface {
	GetBlock(ctx context.Context, blockId starknetData.BlockID) (Block, error)
	TraceBlock(ctx context.Context, blockId starknetData.BlockID) ([]starknet.Trace, error)
	GetStateUpdate(ctx context.Context, blockId starknetData.BlockID) (starknetData.StateUpdate, error)
	GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error)
	TransactionStatus(ctx context.Context, hash string) (storage.Status, error)
	GetClass(ctx context.Context, hash string) (starknetData.Class, error)
	Head(ctx context.Context) (uint64, error)
}

type Block struct {
	Height           uint64
	Time             time.Time
	Version          *string
	Status           storage.Status
	Hash             []byte
	ParentHash       []byte
	NewRoot          []byte
	SequencerAddress []byte

	Transactions []Transaction
	Receipts     []starknet.Receipt
}

type Transaction struct {
	Type      string
	Version   data.Felt
	Hash      data.Felt
	ActualFee data.Felt

	Body any
}
