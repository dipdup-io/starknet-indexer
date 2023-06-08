package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// outputs names
const (
	OutputMessages = "messages"
)

// Stream -
type Stream struct {
	stream  *grpcSDK.Stream[pb.Subscription]
	request *pb.SubscribeRequest
	id      uint64
}

// NewStream -
func NewStream(stream *grpcSDK.Stream[pb.Subscription], request *pb.SubscribeRequest, id uint64) *Stream {
	return &Stream{
		request: request,
		stream:  stream,
		id:      id,
	}
}

// Client -
type Client struct {
	grpc *grpcSDK.Client

	output  *modules.Output
	streams map[uint64]*Stream

	service   pb.IndexerServiceClient
	reconnect chan uint64

	wg *sync.WaitGroup
}

// NewClient -
func NewClient(cfg ClientConfig) *Client {
	return &Client{
		grpc:      grpcSDK.NewClient(cfg.ServerAddress),
		output:    modules.NewOutput(OutputMessages),
		streams:   make(map[uint64]*Stream),
		reconnect: make(chan uint64, 16),
		wg:        new(sync.WaitGroup),
	}
}

// NewClientWithServerAddress -
func NewClientWithServerAddress(address string) *Client {
	return &Client{
		grpc:      grpcSDK.NewClient(address),
		output:    modules.NewOutput(OutputMessages),
		streams:   make(map[uint64]*Stream),
		reconnect: make(chan uint64, 16),
		wg:        new(sync.WaitGroup),
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
	client.grpc.Start(ctx)
	client.service = pb.NewIndexerServiceClient(client.grpc.Connection())

	client.wg.Add(1)
	go client.reconnectThread(ctx)
}

// Connect -
func (client *Client) Connect(ctx context.Context, opts ...grpcSDK.ConnectOption) error {
	return client.grpc.Connect(ctx, opts...)
}

// Close - closes client
func (client *Client) Close() error {
	client.wg.Wait()

	for id, stream := range client.streams {
		unsubscribeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Uint64("id", id).Msg("unsubscribing...")
		if err := stream.stream.Unsubscribe(unsubscribeCtx, id); err != nil {
			log.Err(err).Msg("unsubscribe error")
		}

		if err := stream.stream.Close(); err != nil {
			return err
		}
	}

	if err := client.grpc.Close(); err != nil {
		return err
	}

	close(client.reconnect)
	return nil
}

// Reconnect -
func (client *Client) Reconnect() <-chan uint64 {
	return client.reconnect
}

func (client *Client) reconnectThread(ctx context.Context) {
	defer client.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-client.grpc.Reconnect():
			newStreams := make(map[uint64]*Stream)
			for id, stream := range client.streams {
				if err := stream.stream.Close(); err != nil {
					log.Err(err).Msg("closing stream after reconnect")
				}
				delete(client.streams, id)
				client.reconnect <- id
			}
			client.streams = newStreams
		}
	}
}

func (client *Client) subscribe(ctx context.Context, req *pb.SubscribeRequest) (uint64, *grpcSDK.Stream[pb.Subscription], error) {
	stream, err := client.service.Subscribe(ctx, req)
	if err != nil {
		return 0, nil, err
	}
	grpcStream := grpc.NewStream[pb.Subscription](stream)

	client.wg.Add(1)
	go client.handleMessage(ctx, grpcStream)

	id, err := grpcStream.Subscribe(ctx)
	return id, grpcStream, err
}

// Subscribe -
func (client *Client) Subscribe(ctx context.Context, req *pb.SubscribeRequest) (uint64, error) {
	id, grpcStream, err := client.subscribe(ctx, req)
	if err != nil {
		return 0, err
	}

	client.streams[id] = NewStream(grpcStream, req, id)
	return id, nil
}

func (client *Client) sendToOutput(name string, data any) error {
	output, err := client.Output(name)
	if err != nil {
		return err
	}
	output.Push(data)
	return nil
}

func (client *Client) handleMessage(ctx context.Context, stream *grpcSDK.Stream[pb.Subscription]) {
	defer client.wg.Done()

	for {
		select {
		case <-stream.Context().Done():
			log.Warn().Msg("stream handler was stopped")
			return
		case msg := <-stream.Listen():
			if err := client.sendToOutput(OutputMessages, msg); err != nil {
				log.Err(err).Msg("sending message to output")
			}
		}
	}
}

// Unsubscribe -
func (client *Client) Unsubscribe(ctx context.Context, id uint64) error {
	if _, err := client.service.Unsubscribe(ctx, &generalPB.UnsubscribeRequest{
		Id: id,
	}); err != nil {
		return err
	}

	return nil
}
