package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	"github.com/pkg/errors"
)

// input names
const (
	InputBlocks = "blocks"
)

// Server -
type Server struct {
	*grpcSDK.Server
	pb.UnimplementedIndexerServiceServer

	db postgres.Storage

	input         *modules.Input
	subscriptions *grpcSDK.Subscriptions[*subscriptions.Message, *pb.Subscription]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(
	cfg *grpcSDK.ServerConfig,
	db postgres.Storage,
) (*Server, error) {
	server, err := grpcSDK.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:        server,
		db:            db,
		input:         modules.NewInput(InputBlocks),
		subscriptions: grpcSDK.NewSubscriptions[*subscriptions.Message, *pb.Subscription](),

		wg: new(sync.WaitGroup),
	}, nil
}

// Start -
func (server *Server) Start(ctx context.Context) {
	pb.RegisterIndexerServiceServer(server.Server.Server(), server)

	server.Server.Start(ctx)

	server.wg.Add(1)
	go server.listen(ctx)
}

func (server *Server) listen(ctx context.Context) {
	defer server.wg.Done()

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		case msg, ok := <-server.input.Listen():
			if !ok {
				return
			}
			switch typedMsg := msg.(type) {
			case *storage.Block:
				server.blockHandler(typedMsg)
			case []*storage.TokenBalance:
			}
		}
	}
}

func (module *Server) blockHandler(block *storage.Block) {
	module.subscriptions.NotifyAll(
		subscriptions.NewBlockMessage(block),
		SubscriptionBlock,
	)

	for i := range block.Declare {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeclareMessage(&block.Declare[i]),
			SubscriptionDeclare,
		)

		module.notifyAboutFee(block.Declare[i].Fee)
		module.notifyAboutInternals(block.Declare[i].Internals)
		module.notifyAboutEvents(block.Declare[i].Events)
		module.notifyAboutMessages(block.Declare[i].Messages)
		module.notifyAboutTransfers(block.Declare[i].Transfers)
	}
	for i := range block.Deploy {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeployMessage(&block.Deploy[i]),
			SubscriptionDeploy,
		)

		module.notifyAboutFee(block.Deploy[i].Fee)
		module.notifyAboutInternals(block.Deploy[i].Internals)
		module.notifyAboutEvents(block.Deploy[i].Events)
		module.notifyAboutMessages(block.Deploy[i].Messages)
		module.notifyAboutTransfers(block.Deploy[i].Transfers)
	}
	for i := range block.DeployAccount {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeployAccountMessage(&block.DeployAccount[i]),
			SubscriptionDeployAccount,
		)

		module.notifyAboutFee(block.DeployAccount[i].Fee)
		module.notifyAboutInternals(block.DeployAccount[i].Internals)
		module.notifyAboutEvents(block.DeployAccount[i].Events)
		module.notifyAboutMessages(block.DeployAccount[i].Messages)
		module.notifyAboutTransfers(block.DeployAccount[i].Transfers)
	}
	for i := range block.Invoke {
		module.subscriptions.NotifyAll(
			subscriptions.NewInvokeMessage(&block.Invoke[i]),
			SubscriptionInvoke,
		)

		module.notifyAboutFee(block.Invoke[i].Fee)
		module.notifyAboutInternals(block.Invoke[i].Internals)
		module.notifyAboutEvents(block.Invoke[i].Events)
		module.notifyAboutMessages(block.Invoke[i].Messages)
		module.notifyAboutTransfers(block.Invoke[i].Transfers)
	}
	for i := range block.L1Handler {
		module.subscriptions.NotifyAll(
			subscriptions.NewL1HandlerMessage(&block.L1Handler[i]),
			SubscriptionL1Handler,
		)

		module.notifyAboutFee(block.L1Handler[i].Fee)
		module.notifyAboutInternals(block.L1Handler[i].Internals)
		module.notifyAboutEvents(block.L1Handler[i].Events)
		module.notifyAboutMessages(block.L1Handler[i].Messages)
		module.notifyAboutTransfers(block.L1Handler[i].Transfers)
	}

	for i := range block.StorageDiffs {
		module.subscriptions.NotifyAll(
			subscriptions.NewStorageDiffMessage(&block.StorageDiffs[i]),
			SubscriptionStorageDiff,
		)
	}

	module.subscriptions.NotifyAll(
		subscriptions.NewEndMessage(),
		SubscriptionEnd,
	)
}

func (module *Server) notifyAboutInternals(txs []storage.Internal) {
	for j := range txs {
		module.subscriptions.NotifyAll(
			subscriptions.NewInternalMessage(&txs[j]),
			SubscriptionInternal,
		)
	}
}

func (module *Server) notifyAboutEvents(events []storage.Event) {
	for j := range events {
		module.subscriptions.NotifyAll(
			subscriptions.NewEventMessage(&events[j]),
			SubscriptionEvent,
		)
	}
}

func (module *Server) notifyAboutMessages(msgs []storage.Message) {
	for j := range msgs {
		module.subscriptions.NotifyAll(
			subscriptions.NewStarknetMessage(&msgs[j]),
			SubscriptionMessage,
		)
	}
}

func (module *Server) notifyAboutTransfers(transfers []storage.Transfer) {
	for j := range transfers {
		module.subscriptions.NotifyAll(
			subscriptions.NewTransferMessage(&transfers[j]),
			SubscriptionTransfer,
		)
		module.notifyAboutTokenBalances(transfers[j].TokenBalanceUpdates())
	}
}

func (module *Server) notifyAboutTokenBalances(balances []storage.TokenBalance) {
	for j := range balances {
		module.subscriptions.NotifyAll(
			subscriptions.NewTokenBalanceMessage(&balances[j]),
			SubscriptionTokenBalance,
		)
	}
}

func (module *Server) notifyAboutFee(fee *storage.Fee) {
	if fee == nil {
		return
	}
	module.subscriptions.NotifyAll(
		subscriptions.NewFeeMessage(fee),
		SubscriptionFee,
	)

	module.notifyAboutInternals(fee.Internals)
	module.notifyAboutEvents(fee.Events)
	module.notifyAboutMessages(fee.Messages)
	module.notifyAboutTransfers(fee.Transfers)
}

// Close -
func (module *Server) Close() error {
	module.wg.Wait()

	if err := module.input.Close(); err != nil {
		return err
	}

	return module.Server.Close()
}

// Input -
func (module *Server) Input(name string) (*modules.Input, error) {
	if name != InputBlocks {
		return nil, errors.Wrap(modules.ErrUnknownInput, name)
	}
	return module.input, nil
}

// Output -
func (module *Server) Output(name string) (*modules.Output, error) {
	return nil, errors.Wrap(modules.ErrUnknownOutput, name)
}

// AttachTo -
func (module *Server) AttachTo(name string, input *modules.Input) error {
	return errors.Wrap(modules.ErrUnknownOutput, name)
}

// Name -
func (module *Server) Name() string {
	return "layer1_grpc_server"
}
