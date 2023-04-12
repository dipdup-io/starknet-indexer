package grpc

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// //////////////////////////////////////////////
// ////////////    HANDLERS    //////////////////
// //////////////////////////////////////////////
// Subscribe -
func (module *Server) Subscribe(req *pb.SubscribeRequest, stream pb.IndexerService_SubscribeServer) error {
	subscription, err := subscriptions.NewSubscription(stream.Context(), module.db, req)
	if err != nil {
		return err
	}
	return grpcSDK.DefaultSubscribeOn[*subscriptions.Message, *pb.Subscription](
		stream,
		module.subscriptions,
		subscription,
		func(id uint64) error {
			return module.sync(stream.Context(), id, req, stream)
		},
	)
}

// Unsubscribe -
func (module *Server) Unsubscribe(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
	return grpcSDK.DefaultUnsubscribe(ctx, module.subscriptions, req.Id)
}
