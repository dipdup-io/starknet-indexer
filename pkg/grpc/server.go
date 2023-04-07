package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
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

	input         *modules.Input
	subscriptions *grpcSDK.Subscriptions[*subscriptions.Message, *pb.Subscription]

	wg *sync.WaitGroup
}

// NewServer -
func NewServer(
	cfg *grpcSDK.ServerConfig,
) (*Server, error) {
	server, err := grpcSDK.NewServer(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		Server:        server,
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
			server.blockHandler(msg.(*storage.Block))
		}
	}
}

func (module *Server) blockHandler(block *storage.Block) {
	module.subscriptions.NotifyAll(
		subscriptions.NewBlockMessage(block),
		SubscriptionBlock,
	)

	module.subscriptions.NotifyAll(
		subscriptions.NewEndMessage(),
		SubscriptionEnd,
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
