package parser

import (
	"context"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/generator"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/interfaces"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	v0 "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/version/v0"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/pkg/errors"
)

func createParser(
	version *string,
	resolver resolver.Resolver,
	cache *cache.Cache,
	blocks storage.IBlock,
) (interfaces.Parser, error) {
	if version == nil {
		return v0.New(resolver, cache, blocks), nil
	}

	switch *version {
	case "0.9.1", "0.10.0", "0.10.1", "0.10.2", "0.10.3", "0.11.0", "0.11.0.2", "0.11.1", "0.11.2", "0.12.0":
		return v0.New(resolver, cache, blocks), nil
	default:
		return nil, errors.Errorf("unknown starknet version: %s", *version)
	}
}

// Parse -
func Parse(
	ctx context.Context,
	receiver *receiver.Receiver,
	cache *cache.Cache,
	idGenerator *generator.IdGenerator,
	blocks storage.IBlock,
	proxies storage.IProxy,
	result receiver.Result,
) (parserData.Result, error) {
	block := storage.Block{
		ID:               result.Block.BlockNumber + 1,
		Height:           result.Block.BlockNumber,
		Time:             time.Unix(result.Block.Timestamp, 0).UTC(),
		Hash:             data.Felt(result.Block.BlockHash).Bytes(),
		ParentHash:       data.Felt(result.Block.ParentHash).Bytes(),
		NewRoot:          encoding.MustDecodeHex(result.Block.NewRoot),
		SequencerAddress: encoding.MustDecodeHex(result.Block.SequencerAddress),
		Version:          result.Block.StarknetVersion,
		Status:           storage.NewStatus(result.Block.Status),
		TxCount:          len(result.Block.Transactions),

		Invoke:        make([]storage.Invoke, 0),
		Declare:       make([]storage.Declare, 0),
		Deploy:        make([]storage.Deploy, 0),
		DeployAccount: make([]storage.DeployAccount, 0),
		L1Handler:     make([]storage.L1Handler, 0),
	}

	blockCtx := parserData.NewBlockContext(block)

	resolver := resolver.NewResolver(receiver, cache, idGenerator, blocks, proxies, blockCtx)

	if err := resolver.ResolveStateUpdates(ctx, &block, result.StateUpdate); err != nil {
		return parserData.Result{}, errors.Wrap(err, "state update parsing")
	}

	p, err := createParser(block.Version, resolver, cache, blocks)
	if err != nil {
		return parserData.Result{}, errors.Wrap(err, "createParser")
	}

	for i := range result.Block.Transactions {
		switch typed := result.Block.Transactions[i].Body.(type) {
		case *starknetData.Invoke:
			var (
				invoke storage.Invoke
				fee    *storage.Fee
				err    error
			)
			switch result.Block.Transactions[i].Version {
			case starknetData.Version0:
				invoke, fee, err = p.ParseInvokeV0(ctx, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			case starknetData.Version1:
				invoke, fee, err = p.ParseInvokeV1(ctx, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			default:
				return parserData.Result{}, errors.Errorf("unknown invoke version: %s", result.Block.Transactions[i].Version)
			}
			if err != nil {
				return parserData.Result{}, errors.Wrapf(err, "%s invoke version=%s", result.Block.Transactions[i].TransactionHash, result.Block.Transactions[i].Version)
			}
			invoke.Position = i
			block.Invoke = append(block.Invoke, invoke)
			if fee != nil {
				block.Fee = append(block.Fee, *fee)
			}
		case *starknetData.Declare:
			tx, fee, err := p.ParseDeclare(ctx, result.Block.Transactions[i].Version, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			if err != nil {
				return parserData.Result{}, errors.Wrapf(err, "%s declare", result.Block.Transactions[i].TransactionHash)
			}
			tx.Position = i
			block.Declare = append(block.Declare, tx)
			if fee != nil {
				block.Fee = append(block.Fee, *fee)
			}
		case *starknetData.Deploy:
			tx, fee, err := p.ParseDeploy(ctx, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			if err != nil {
				return parserData.Result{}, errors.Wrapf(err, "%s deploy", result.Block.Transactions[i].TransactionHash)
			}
			tx.Position = i
			block.Deploy = append(block.Deploy, tx)
			if fee != nil {
				block.Fee = append(block.Fee, *fee)
			}
		case *starknetData.DeployAccount:
			tx, fee, err := p.ParseDeployAccount(ctx, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			if err != nil {
				return parserData.Result{}, errors.Wrapf(err, "%s deploy account", result.Block.Transactions[i].TransactionHash)
			}
			tx.Position = i
			block.DeployAccount = append(block.DeployAccount, tx)
			if fee != nil {
				block.Fee = append(block.Fee, *fee)
			}
		case *starknetData.L1Handler:
			tx, fee, err := p.ParseL1Handler(ctx, typed, block, result.Trace.Traces[i], result.Block.Receipts[i])
			if err != nil {
				return parserData.Result{}, errors.Wrapf(err, "%s l1 handler", result.Block.Transactions[i].TransactionHash)
			}
			tx.Position = i
			block.L1Handler = append(block.L1Handler, tx)
			if fee != nil {
				block.Fee = append(block.Fee, *fee)
			}
		default:
			return parserData.Result{}, errors.Errorf("unknown transaction type: %s", result.Block.Transactions[i].Type)
		}
	}

	block.InvokeCount = len(block.Invoke)
	block.DeclareCount = len(block.Declare)
	block.DeployCount = len(block.Deploy)
	block.DeployAccountCount = len(block.DeployAccount)
	block.L1HandlerCount = len(block.L1Handler)

	return parserData.Result{
		Block:   block,
		Context: blockCtx,
	}, nil
}
