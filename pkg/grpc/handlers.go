package grpc

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	grpcSDK "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

// //////////////////////////////////////////////
// ////////////    HANDLERS    //////////////////
// //////////////////////////////////////////////

// Subscribe -
func (module *Server) Subscribe(req *pb.SubscribeRequest, stream pb.IndexerService_SubscribeServer) error {
	module.log.Info().Msg("subscribe request")
	subscription, err := subscriptions.NewSubscription(stream.Context(), module.db, req)
	if err != nil {
		return err
	}
	var height uint64
	return grpcSDK.DefaultSubscribeOn[*subscriptions.Message, *pb.Subscription](
		stream,
		module.subscriptions,
		subscription,
		func(id uint64) error {
			height, err = module.sync(stream.Context(), id, req, stream)
			return err
		},
		func(id uint64) error {
			return stream.Send(&pb.Subscription{
				Response: &generalPB.SubscribeResponse{
					Id: id,
				},
				EndOfBlock: &pb.EndOfBlock{
					Height: height,
				},
			})
		},
	)
}

// Unsubscribe -
func (module *Server) Unsubscribe(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
	module.log.Info().Msg("unsubscribe request")
	return grpcSDK.DefaultUnsubscribe(ctx, module.subscriptions, req.Id)
}

// JSONSchemaForClass -
func (module *Server) JSONSchemaForClass(ctx context.Context, req *pb.Bytes) (*pb.Bytes, error) {
	if req == nil {
		return nil, errors.Errorf("empty class hash")
	}

	hash := req.GetData()
	module.log.Info().Hex("hash", hash).Msg("json schema for class request")

	if !starknet.HashValidator(hash) {
		return nil, errors.Errorf("invalid starknet hash (length must be 32): %x", hash)
	}
	class, err := module.db.Class.GetByHash(ctx, hash)
	if err != nil {
		return nil, errors.Wrapf(err, "receiving class error %x", hash)
	}

	if len(class.Abi) == 0 {
		return nil, errors.Errorf("empty abi for class %x", class.Hash)
	}

	var a abi.Abi
	if err := json.Unmarshal(class.Abi, &a); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal abi %x", hash)
	}

	schema := a.JsonSchema()
	b, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal json schema")
	}
	return &pb.Bytes{
		Data: b,
	}, nil
}

// JSONSchemaForContract -
func (module *Server) JSONSchemaForContract(ctx context.Context, req *pb.Bytes) (*pb.Bytes, error) {
	if req == nil {
		return nil, errors.Errorf("empty contract hash")
	}

	hash := req.GetData()
	module.log.Info().Hex("hash", hash).Msg("json schema for contract request")

	if !starknet.HashValidator(hash) {
		return nil, errors.Errorf("invalid starknet hash (length must be 32): %x", hash)
	}
	contract, err := module.db.Address.GetByHash(ctx, hash)
	if err != nil {
		return nil, errors.Wrapf(err, "receiving contract error %x", hash)
	}
	if contract.ClassID == nil {
		return nil, errors.Wrapf(err, "unknown class for contract %x", hash)
	}
	class, err := module.db.Class.GetByID(ctx, *contract.ClassID)
	if err != nil {
		return nil, errors.Wrapf(err, "receiving class error for contract %x", hash)
	}

	var a abi.Abi
	if err := json.Unmarshal(class.Abi, &a); err != nil {
		return nil, errors.Wrapf(err, "can't unmarshal abi %x", hash)
	}

	schema := a.JsonSchema()
	b, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal json schema")
	}
	return &pb.Bytes{
		Data: b,
	}, nil
}

// GetProxy -
func (module *Server) GetProxy(ctx context.Context, req *pb.ProxyRequest) (*pb.Proxy, error) {
	if req == nil {
		return nil, errors.Errorf("empty proxy request")
	}

	hash := req.GetHash().GetData()
	selector := req.GetSelector().GetData()
	module.log.Info().Hex("hash", hash).Msg("get proxy request")

	if !starknet.HashValidator(hash) {
		return nil, errors.Errorf("invalid starknet hash (length must be 32): %x", hash)
	}
	proxy, err := module.db.Proxy.GetByHash(ctx, hash, selector)
	if err != nil {
		return nil, errors.Wrapf(err, "receiving proxy error %x", hash)
	}
	return &pb.Proxy{
		Id:   proxy.EntityID,
		Hash: proxy.EntityHash,
		Type: uint32(proxy.EntityType),
	}, nil
}
