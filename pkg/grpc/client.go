package grpc

import (
	"context"
	"time"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/rs/zerolog/log"
)

const (
	OutputMessages = "messages"
	ModuleName     = "layer1_grpc_client"
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
	modules.BaseModule
	grpc    *grpcSDK.Client
	streams map[uint64]*Stream

	service   pb.IndexerServiceClient
	reconnect chan uint64
}

// NewClient -
func NewClient(cfg ClientConfig) *Client {
	client := &Client{
		BaseModule: modules.New(ModuleName),
		grpc:       grpcSDK.NewClient(cfg.ServerAddress),
		streams:    make(map[uint64]*Stream),
		reconnect:  make(chan uint64, 16),
	}
	client.CreateOutput(OutputMessages)
	return client
}

// NewClientWithServerAddress -
func NewClientWithServerAddress(address string) *Client {
	client := &Client{
		BaseModule: modules.New(ModuleName),
		grpc:       grpcSDK.NewClient(address),
		streams:    make(map[uint64]*Stream),
		reconnect:  make(chan uint64, 16),
	}
	client.CreateOutput(OutputMessages)
	return client
}

// Start -
func (client *Client) Start(ctx context.Context) {
	client.grpc.Start(ctx)
	client.service = pb.NewIndexerServiceClient(client.grpc.Connection())
	client.G.GoCtx(ctx, client.reconnectThread)
}

// Connect -
func (client *Client) Connect(ctx context.Context, opts ...grpcSDK.ConnectOption) error {
	return client.grpc.Connect(ctx, opts...)
}

// Close - closes client
func (client *Client) Close() error {
	client.G.Wait()

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

	client.G.GoCtx(ctx, func(ctx context.Context) {
		client.handleMessage(ctx, grpcStream)
	})

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

// JsonSchemaForClass -
func (client *Client) JsonSchemaForClass(ctx context.Context, req *pb.Bytes) (*pb.Bytes, error) {
	return client.service.JSONSchemaForClass(ctx, req)
}

// JsonSchemaForContract -
func (client *Client) JsonSchemaForContract(ctx context.Context, req *pb.Bytes) (*pb.Bytes, error) {
	return client.service.JSONSchemaForContract(ctx, req)
}

// GetProxy -
func (client *Client) GetProxy(ctx context.Context, hash, selector []byte) (*pb.Proxy, error) {
	return client.service.GetProxy(ctx, &pb.ProxyRequest{
		Hash: &pb.Bytes{
			Data: hash,
		},
		Selector: &pb.Bytes{
			Data: selector,
		},
	})
}
