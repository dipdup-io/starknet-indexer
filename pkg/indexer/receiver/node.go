package receiver

import (
	"context"
	"os"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/rpc"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
)

type Node struct {
	api starknet.API
}

func NewNode(cfg config.DataSource) *Node {
	apiKey := os.Getenv("NODE_APIKEY")
	headerName := os.Getenv("NODE_HEADER_APIKEY")
	return &Node{
		api: starknet.NewAPI(
			cfg.URL,
			starknet.WithRateLimit(cfg.RequestsPerSecond),
			starknet.WithApiKey(headerName, apiKey),
		),
	}
}

func (n *Node) GetBlock(ctx context.Context, blockId starknetData.BlockID) (block Block, err error) {
	response, err := n.api.GetBlockWithReceipts(ctx, blockId)
	if err != nil {
		return
	}

	block.Height = response.Result.BlockNumber
	block.Time = time.Unix(response.Result.Timestamp, 0).UTC()
	block.Hash = data.Felt(response.Result.BlockHash).Bytes()
	block.ParentHash = data.Felt(response.Result.ParentHash).Bytes()
	block.NewRoot = encoding.MustDecodeHex(response.Result.NewRoot)
	block.SequencerAddress = encoding.MustDecodeHex(response.Result.SequencerAddress)
	block.Version = response.Result.Version
	block.Status = storage.NewStatus(response.Result.Status)
	block.Transactions = make([]Transaction, len(response.Result.Transactions))

	for i := range response.Result.Transactions {
		block.Transactions[i].Hash = response.Result.Transactions[i].Transaction.TransactionHash
		block.Transactions[i].Type = response.Result.Transactions[i].Transaction.Type
		block.Transactions[i].Version = response.Result.Transactions[i].Transaction.Version
		block.Transactions[i].Body = response.Result.Transactions[i].Transaction.Body
		block.Transactions[i].ActualFee = response.Result.Transactions[i].Receipt.ActualFee.Amount

		switch block.Transactions[i].Type {
		case starknetData.TransactionTypeDeploy:
			if deploy, ok := block.Transactions[i].Body.(*starknetData.Deploy); ok {
				deploy.ContractAddress = starknetData.Felt(response.Result.Transactions[i].Receipt.ContractAddress)
			} else {
				return block, errors.Errorf("invalid invoke transaction type: expected Deploy (non-pointer)")
			}
		case starknetData.TransactionTypeDeployAccount:
			if deploy, ok := block.Transactions[i].Body.(*starknetData.DeployAccount); ok {
				deploy.ContractAddress = starknetData.Felt(response.Result.Transactions[i].Receipt.ContractAddress)
			} else {
				return block, errors.Errorf("invalid invoke transaction type: expected DeployAccount (non-pointer)")
			}
		default:
			continue
		}
	}

	return
}

func (n *Node) TraceBlock(ctx context.Context, block starknetData.BlockID) (traces []sequencer.Trace, err error) {
	response, err := n.api.Trace(ctx, block)
	if err != nil {
		return
	}

	traces = make([]sequencer.Trace, len(response.Result))
	for i := range response.Result {
		if inv := response.Result[i].TraceRoot.ExecuteInvocation; inv != nil {
			if inv.RevertReason != "" {
				traces[i].RevertedError = inv.RevertReason
			} else {
				traces[i].FunctionInvocation = makeSeqInvocationFromNodeCall(inv)
			}
		}

		if inv := response.Result[i].TraceRoot.ConstructorInvocation; inv != nil {
			if inv.RevertReason != "" {
				traces[i].RevertedError = inv.RevertReason
			} else {
				traces[i].FunctionInvocation = makeSeqInvocationFromNodeCall(inv)
			}
		}

		traces[i].ValidateInvocation = makeSeqInvocationFromNodeCall(response.Result[i].TraceRoot.ValidateInvocation)
		traces[i].FeeTransferInvocation = makeSeqInvocationFromNodeCall(response.Result[i].TraceRoot.FeeTransferInvocation)
		traces[i].TransactionHash = response.Result[i].TransactionHash
	}

	return
}

func makeSeqInvocationFromNodeCall(call *starknet.Call) *sequencer.Invocation {
	if call == nil {
		return nil
	}

	inv := &sequencer.Invocation{
		CallerAddress:   call.CallerAddress,
		ContractAddress: call.ContractAddress,
		Calldata:        call.Calldata,
		CallType:        call.CallType,
		ClassHash:       call.ClassHash,
		Selector:        call.EntryPointSelector,
		EntrypointType:  call.EntryPointType,
		Result:          call.Result,
		Events:          call.Events,
		Messages:        call.Messages,
		InternalCalls:   make([]sequencer.Invocation, len(call.Calls)),
	}

	for i := range call.Calls {
		internalCall := makeSeqInvocationFromNodeCall(&call.Calls[i])
		inv.InternalCalls[i] = *internalCall
	}

	return inv
}

func (n *Node) GetStateUpdate(ctx context.Context, block starknetData.BlockID) (starknetData.StateUpdate, error) {
	response, err := n.api.GetStateUpdate(ctx, block)
	if err != nil {
		return starknetData.StateUpdate{}, err
	}
	return response.Result.ToStateUpdate(), nil
}

func (n *Node) GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error) {
	response, err := n.api.GetBlockWithTxHashes(ctx, starknetData.BlockID{Number: &height})
	if err != nil {
		return storage.StatusUnknown, err
	}
	return storage.NewStatus(response.Result.Status), nil
}

func (n *Node) TransactionStatus(ctx context.Context, hash string) (storage.Status, error) {
	response, err := n.api.GetTransactionStatus(ctx, hash)
	if err != nil {
		return storage.StatusUnknown, err
	}

	return storage.NewStatus(response.Result.Finality), nil
}

func (n *Node) GetClass(ctx context.Context, hash string) (starknetData.Class, error) {
	blockId := starknetData.BlockID{
		String: starknetData.Latest,
	}

	response, err := n.api.GetClass(ctx, blockId, hash)
	if err != nil {
		return starknetData.Class{}, err
	}
	return response.Result, nil
}

func (n *Node) Head(ctx context.Context) (uint64, error) {
	response, err := n.api.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return response.Result, nil
}
