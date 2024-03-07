package interfaces

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
)

// Parser -
type Parser interface {
	ParseDeclare(ctx context.Context, version starknetData.Felt, raw *starknetData.Declare, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Declare, *storage.Fee, error)
	ParseDeployAccount(ctx context.Context, raw *starknetData.DeployAccount, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.DeployAccount, *storage.Fee, error)
	ParseDeploy(ctx context.Context, raw *starknetData.Deploy, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Deploy, *storage.Fee, error)
	ParseInvokeV0(ctx context.Context, raw *starknetData.Invoke, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Invoke, *storage.Fee, error)
	ParseInvokeV1(ctx context.Context, raw *starknetData.Invoke, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Invoke, *storage.Fee, error)
	ParseInvokeV3(ctx context.Context, raw *starknetData.Invoke, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.Invoke, *storage.Fee, error)
	ParseL1Handler(ctx context.Context, raw *starknetData.L1Handler, block storage.Block, receiverTx receiver.Transaction, trace sequencer.Trace) (storage.L1Handler, *storage.Fee, error)
}

// InternalTxParser -
type InternalTxParser interface {
	Parse(ctx context.Context, txCtx data.TxContext, internal sequencer.Invocation) (storage.Internal, error)
}

// EventParser -
type EventParser interface {
	Parse(ctx context.Context, txCtx data.TxContext, contractAbi abi.Abi, event starknetData.Event) (storage.Event, error)
}

// MessageParser -
type MessageParser interface {
	Parse(ctx context.Context, txCtx data.TxContext, msg starknetData.Message) (storage.Message, error)
}

// FeeParser -
type FeeParser interface {
	ParseInvocation(ctx context.Context, txCtx data.TxContext, feeInvocation sequencer.Invocation) (*storage.Fee, error)
	ParseActualFee(ctx context.Context, txCtx data.TxContext, actualFee starknetData.Felt) (*storage.Transfer, error)
}

// TransferParser -
type TransferParser interface {
	ParseEvents(ctx context.Context, txCtx data.TxContext, contract storage.Address, events []storage.Event) ([]storage.Transfer, error)
	ParseCalldata(ctx context.Context, txCtx data.TxContext, entrypoint string, calldata map[string]any) ([]storage.Transfer, error)
}

// ProxyUpgrader -
type ProxyUpgrader interface {
	Parse(ctx context.Context, txCtx data.TxContext, contract storage.Address, events []storage.Event, entrypoint string, data map[string]any) error
}
