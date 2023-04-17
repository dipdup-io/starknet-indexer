package grpc

import (
	"context"
	"sync"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
)

// outputs names
const (
	OutputMessages = "messages"
)

// Client -
type Client struct {
	*grpcSDK.Client

	output *modules.Output

	client pb.IndexerServiceClient

	wg *sync.WaitGroup
}

// NewClient -
func NewClient(cfg ClientConfig) *Client {
	return &Client{
		Client: grpcSDK.NewClient(cfg.ServerAddress),
		output: modules.NewOutput(OutputMessages),
		wg:     new(sync.WaitGroup),
	}
}

// NewClientWithServerAddress -
func NewClientWithServerAddress(address string) *Client {
	return &Client{
		Client: grpcSDK.NewClient(address),
		output: modules.NewOutput(OutputMessages),
		wg:     new(sync.WaitGroup),
	}
}

// Name -
func (client *Client) Name() string {
	return "layer1_grpc_client"
}

// Input -
func (client *Client) Input(name string) (*modules.Input, error) {
	return nil, errors.Wrap(modules.ErrUnknownInput, name)
}

// Output -
func (client *Client) Output(name string) (*modules.Output, error) {
	if name != OutputMessages {
		return nil, errors.Wrap(modules.ErrUnknownOutput, name)
	}
	return client.output, nil
}

// AttachTo -
func (client *Client) AttachTo(name string, input *modules.Input) error {
	output, err := client.Output(name)
	if err != nil {
		return err
	}
	output.Attach(input)
	return nil
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.client = pb.NewIndexerServiceClient(client.Connection())
}

// Subscribe -
func (client *Client) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (uint64, error) {
	stream, err := client.client.Subscribe(ctx, req)
	if err != nil {
		return 0, err
	}

	return grpc.Subscribe[*pb.Subscription](
		stream,
		client.handleMessage,
		client.wg,
	)
}

func (client *Client) sendToOutput(name string, data any) error {
	output, err := client.Output(name)
	if err != nil {
		return err
	}
	output.Push(data)
	return nil
}

func (client *Client) handleMessage(ctx context.Context, data *pb.Subscription, id uint64) error {
	return client.sendToOutput(OutputMessages, data)
}

// Unsubscribe -
func (client *Client) Unsubscribe(ctx context.Context, id uint64) error {
	if _, err := client.client.Unsubscribe(ctx, &generalPB.UnsubscribeRequest{
		Id: id,
	}); err != nil {
		return err
	}

	return nil
}
