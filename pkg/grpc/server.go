package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	log           zerolog.Logger

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
		log:           log.With().Str("module", "grpc_server").Logger(),

		wg: new(sync.WaitGroup),
	}, nil
}

// Start -
func (module *Server) Start(ctx context.Context) {
	pb.RegisterIndexerServiceServer(module.Server.Server(), module)

	module.Server.Start(ctx)

	module.wg.Add(1)
	go module.listen(ctx)
}

func (module *Server) listen(ctx context.Context) {
	defer module.wg.Done()

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		case msg, ok := <-module.input.Listen():
			if !ok {
				return
			}
			if message, ok := msg.(*indexer.IndexerMessage); ok {
				module.blockHandler(ctx, message)
			} else {
				module.log.Warn().Msgf("unknown message type: %T", msg)
			}
		}
	}
}

func (module *Server) blockHandler(ctx context.Context, message *indexer.IndexerMessage) {
	for _, address := range message.Addresses {
		module.notifyAboutAddress(address)
	}

	module.subscriptions.NotifyAll(
		subscriptions.NewBlockMessage(message.Block),
		SubscriptionBlock,
	)

	for _, token := range message.Tokens {
		if err := module.notifyAboutToken(ctx, token); err != nil {
			log.Err(err).Msg("can't notify about token")
		}
	}

	for i := range message.Block.Declare {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeclareMessage(&message.Block.Declare[i]),
			SubscriptionDeclare,
		)

		module.notifyAboutFee(message.Block.Declare[i].Fee)
		module.notifyAboutInternals(message.Block.Declare[i].Internals)
		module.notifyAboutEvents(message.Block.Declare[i].Events)
		module.notifyAboutMessages(message.Block.Declare[i].Messages)
		module.notifyAboutTransfers(message.Block.Declare[i].Transfers)
	}
	for i := range message.Block.Deploy {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeployMessage(&message.Block.Deploy[i]),
			SubscriptionDeploy,
		)

		module.notifyAboutFee(message.Block.Deploy[i].Fee)
		module.notifyAboutInternals(message.Block.Deploy[i].Internals)
		module.notifyAboutEvents(message.Block.Deploy[i].Events)
		module.notifyAboutMessages(message.Block.Deploy[i].Messages)
		module.notifyAboutTransfers(message.Block.Deploy[i].Transfers)

	}
	for i := range message.Block.DeployAccount {
		module.subscriptions.NotifyAll(
			subscriptions.NewDeployAccountMessage(&message.Block.DeployAccount[i]),
			SubscriptionDeployAccount,
		)

		module.notifyAboutFee(message.Block.DeployAccount[i].Fee)
		module.notifyAboutInternals(message.Block.DeployAccount[i].Internals)
		module.notifyAboutEvents(message.Block.DeployAccount[i].Events)
		module.notifyAboutMessages(message.Block.DeployAccount[i].Messages)
		module.notifyAboutTransfers(message.Block.DeployAccount[i].Transfers)
	}
	for i := range message.Block.Invoke {
		module.subscriptions.NotifyAll(
			subscriptions.NewInvokeMessage(&message.Block.Invoke[i]),
			SubscriptionInvoke,
		)

		module.notifyAboutFee(message.Block.Invoke[i].Fee)
		module.notifyAboutInternals(message.Block.Invoke[i].Internals)
		module.notifyAboutEvents(message.Block.Invoke[i].Events)
		module.notifyAboutMessages(message.Block.Invoke[i].Messages)
		module.notifyAboutTransfers(message.Block.Invoke[i].Transfers)
	}
	for i := range message.Block.L1Handler {
		module.subscriptions.NotifyAll(
			subscriptions.NewL1HandlerMessage(&message.Block.L1Handler[i]),
			SubscriptionL1Handler,
		)

		module.notifyAboutFee(message.Block.L1Handler[i].Fee)
		module.notifyAboutInternals(message.Block.L1Handler[i].Internals)
		module.notifyAboutEvents(message.Block.L1Handler[i].Events)
		module.notifyAboutMessages(message.Block.L1Handler[i].Messages)
		module.notifyAboutTransfers(message.Block.L1Handler[i].Transfers)
	}

	for i := range message.Block.StorageDiffs {
		module.subscriptions.NotifyAll(
			subscriptions.NewStorageDiffMessage(&message.Block.StorageDiffs[i]),
			SubscriptionStorageDiff,
		)
	}

	module.subscriptions.NotifyAll(
		subscriptions.NewEndMessage(message.Block),
		SubscriptionEnd,
	)
}

func (module *Server) notifyAboutInternals(txs []storage.Internal) {
	for j := range txs {
		module.subscriptions.NotifyAll(
			subscriptions.NewInternalMessage(&txs[j]),
			SubscriptionInternal,
		)
		module.notifyAboutInternals(txs[j].Internals)
		module.notifyAboutEvents(txs[j].Events)
		module.notifyAboutMessages(txs[j].Messages)
		module.notifyAboutTransfers(txs[j].Transfers)
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

func (module *Server) notifyAboutToken(ctx context.Context, token *storage.Token) error {
	if token == nil {
		return nil
	}

	contract, err := module.db.Address.GetByID(ctx, token.ContractId)
	if err != nil {
		return err
	}
	token.Contract = *contract

	module.subscriptions.NotifyAll(
		subscriptions.NewTokenMessage(token),
		SubscriptionToken,
	)

	return nil
}

func (module *Server) notifyAboutAddress(address *storage.Address) {
	if address == nil {
		return
	}

	module.subscriptions.NotifyAll(
		subscriptions.NewAddressMessage(address),
		SubscriptionAddress,
	)
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
