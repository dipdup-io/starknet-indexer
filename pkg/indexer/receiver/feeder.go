package receiver

import (
	"context"
	"time"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
)

type Feeder struct {
	api starknet.API
}

func NewFeeder(cfg config.DataSource) *Feeder {
	opts := make([]starknet.ApiOption, 0)
	if cfg.RequestsPerSecond > 0 {
		opts = append(opts, starknet.WithRateLimit(cfg.RequestsPerSecond))
	}

	return &Feeder{
		api: starknet.NewAPI("", cfg.URL, opts...),
	}
}

func (f *Feeder) GetBlock(ctx context.Context, blockId starknetData.BlockID) (block Block, err error) {
	response, err := f.api.GetBlock(ctx, blockId, false)
	if err != nil {
		return block, err
	}

	block.Height = response.BlockNumber
	block.Time = time.Unix(response.Timestamp, 0).UTC()
	block.Hash = starknetData.Felt(response.BlockHash).Bytes()
	block.ParentHash = starknetData.Felt(response.ParentHash).Bytes()
	block.NewRoot = encoding.MustDecodeHex(response.NewRoot)
	block.SequencerAddress = encoding.MustDecodeHex(response.SequencerAddress)
	block.Version = response.StarknetVersion
	block.Status = storage.NewStatus(response.Status)
	block.Receipts = response.Receipts

	if len(response.Transactions) != len(response.Receipts) {
		return block, errors.Errorf("length arrays of txs and receipts are differ")
	}

	block.Transactions = make([]Transaction, len(response.Transactions))

	for i := range response.Transactions {
		block.Transactions[i].Hash = response.Transactions[i].TransactionHash
		block.Transactions[i].Type = response.Transactions[i].Type
		block.Transactions[i].Version = response.Transactions[i].Version
		block.Transactions[i].Body = response.Transactions[i].Body
		block.Transactions[i].ActualFee = response.Receipts[i].ActualFee
	}

	return
}

func (f *Feeder) TraceBlock(ctx context.Context, blockId starknetData.BlockID) (traces []starknet.Trace, err error) {
	response, err := f.api.TraceBlock(ctx, blockId)
	if err != nil {
		return
	}
	return response.Traces, nil
}

func (f *Feeder) GetStateUpdate(ctx context.Context, blockId starknetData.BlockID) (response starknetData.StateUpdate, err error) {
	return f.api.GetStateUpdate(ctx, blockId)
}

func (f *Feeder) GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error) {
	response, err := f.api.GetBlock(ctx, starknetData.BlockID{Number: &height}, false)
	if err != nil {
		return storage.StatusUnknown, err
	}
	return storage.NewStatus(response.Status), nil
}

func (f *Feeder) TransactionStatus(ctx context.Context, hash string) (storage.Status, error) {
	response, err := f.api.GetTransactionStatus(ctx, hash)
	if err != nil {
		return storage.StatusUnknown, err
	}

	return storage.NewStatus(response.Status), nil
}

func (f *Feeder) GetClass(ctx context.Context, hash string) (starknetData.Class, error) {
	blockId := starknetData.BlockID{
		String: starknetData.Latest,
	}

	return f.api.GetClassByHash(ctx, blockId, hash)
}

func (f *Feeder) Head(ctx context.Context) (uint64, error) {
	response, err := f.api.GetBlock(ctx, starknetData.BlockID{
		String: starknetData.Latest,
	}, true)
	if err != nil {
		return 0, err
	}
	return response.BlockNumber, nil
}
