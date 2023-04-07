package grpc

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
)

// SubscriptionBlock -
func SubscriptionBlock(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Block: Block(msg.Block),
	}
}

// Block -
func Block(block *storage.Block) *pb.Block {
	pbBlock := &pb.Block{
		Id:                 block.ID,
		Height:             block.Height,
		Time:               uint64(block.Time.Unix()),
		TxCount:            uint64(block.TxCount),
		InvokesCount:       uint64(block.InvokeCount),
		DeclaresCount:      uint64(block.DeclareCount),
		DeploysCount:       uint64(block.DeployCount),
		DeployAccountCount: uint64(block.DeployAccountCount),
		L1HandlersCount:    uint64(block.L1HandlerCount),
		StorageDiffsCount:  uint64(block.StorageDiffCount),
		Status:             uint64(block.Status),
		Hash:               block.Hash,
		ParentHash:         block.ParentHash,
		NewRoot:            block.NewRoot,
		SequencerAddress:   block.SequencerAddress,
	}

	if block.Version != nil {
		pbBlock.Version = *block.Version
	}

	return pbBlock
}

// SubscriptionEnd -
func SubscriptionEnd(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		EndOfBlock: true,
	}
}
