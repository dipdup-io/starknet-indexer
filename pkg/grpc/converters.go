package grpc

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"github.com/goccy/go-json"
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

// SubscriptionDeclare -
func SubscriptionDeclare(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Declare: Declare(msg.Declare),
	}
}

// Declare -
func Declare(model *storage.Declare) *pb.Declare {
	pbDeclare := &pb.Declare{
		Id:       model.ID,
		Height:   model.Height,
		Time:     uint64(model.Time.Unix()),
		Version:  model.Version,
		Position: uint64(model.Position),
		Sender:   model.Sender.Hash,
		Contract: model.Contract.Hash,
		Class:    model.Class.Hash,
		Status:   uint64(model.Status),
		Hash:     model.Hash,
		MaxFee:   model.MaxFee.String(),
		Nonce:    model.Nonce.String(),
	}
	if model.ContractID != nil {
		pbDeclare.ContractId = *model.ContractID
	}
	if model.SenderID != nil {
		pbDeclare.SenderId = *model.SenderID
	}
	return pbDeclare
}

// SubscriptionDeploy -
func SubscriptionDeploy(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Deploy: Deploy(msg.Deploy),
	}
}

// Deploy -
func Deploy(model *storage.Deploy) *pb.Deploy {
	pbModel := &pb.Deploy{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Position:   uint64(model.Position),
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		Class:      model.Class.Hash,
		Status:     uint64(model.Status),
		Hash:       model.Hash,
		Salt:       model.ContractAddressSalt,
		Calldata:   model.ConstructorCalldata,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}

	return pbModel
}

// SubscriptionDeployAccount -
func SubscriptionDeployAccount(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		DeployAccount: DeployAccount(msg.DeployAccount),
	}
}

// DeployAccount -
func DeployAccount(model *storage.DeployAccount) *pb.DeployAccount {
	pbModel := &pb.DeployAccount{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Position:   uint64(model.Position),
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		Class:      model.Class.Hash,
		Status:     uint64(model.Status),
		Hash:       model.Hash,
		MaxFee:     model.MaxFee.String(),
		Nonce:      model.Nonce.String(),
		Salt:       model.ContractAddressSalt,
		Calldata:   model.ConstructorCalldata,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}

	return pbModel
}

// SubscriptionEvent -
func SubscriptionEvent(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Event: Event(msg.Event),
	}
}

// Event -
func Event(model *storage.Event) *pb.Event {
	pbModel := &pb.Event{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Order:      model.Order,
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		From:       model.From.Hash,
		FromId:     model.FromID,
		Keys:       model.Keys,
		Data:       model.Data,
		Name:       model.Name,
	}

	if model.ParsedData != nil {
		parsed, err := json.Marshal(model.ParsedData)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedData = parsed
	}

	return pbModel
}

// SubscriptionFee -
func SubscriptionFee(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Fee: Fee(msg.Fee),
	}
}

// Fee -
func Fee(model *storage.Fee) *pb.Fee {
	pbModel := &pb.Fee{
		Id:             model.ID,
		Height:         model.Height,
		Time:           uint64(model.Time.Unix()),
		Contract:       model.Contract.Hash,
		ContractId:     model.ContractID,
		Caller:         model.Caller.Hash,
		CallerId:       model.CallerID,
		Class:          model.Class.Hash,
		Selector:       model.Selector,
		EntrypointType: uint64(model.EntrypointType),
		CallType:       uint64(model.CallType),
		Calldata:       model.Calldata,
		Result:         model.Result,
		Entrypoint:     model.Entrypoint,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}

	return pbModel
}

// SubscriptionInternal -
func SubscriptionInternal(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Internal: Internal(msg.Internal),
	}
}

// Internal -
func Internal(model *storage.Internal) *pb.Internal {
	pbModel := &pb.Internal{
		Id:             model.ID,
		Height:         model.Height,
		Time:           uint64(model.Time.Unix()),
		Status:         uint64(model.Status),
		Hash:           model.Hash,
		Contract:       model.Contract.Hash,
		ContractId:     model.ContractID,
		Caller:         model.Caller.Hash,
		CallerId:       model.CallerID,
		Class:          model.Class.Hash,
		Selector:       model.Selector,
		EntrypointType: uint64(model.EntrypointType),
		CallType:       uint64(model.CallType),
		Calldata:       model.Calldata,
		Result:         model.Result,
		Entrypoint:     model.Entrypoint,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}
	if model.ParsedResult != nil {
		parsed, err := json.Marshal(model.ParsedResult)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedResult = parsed
	}
	return pbModel
}

// SubscriptionInvoke -
func SubscriptionInvoke(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Invoke: Invoke(msg.Invoke),
	}
}

// Invoke -
func Invoke(model *storage.Invoke) *pb.Invoke {
	pbModel := &pb.Invoke{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Status:     uint64(model.Status),
		Hash:       model.Hash,
		Position:   uint64(model.Position),
		Version:    model.Version,
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		Selector:   model.EntrypointSelector,
		Calldata:   model.CallData,
		MaxFee:     model.MaxFee.String(),
		Nonce:      model.Nonce.String(),
		Entrypoint: model.Entrypoint,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}
	return pbModel
}

// SubscriptionL1Handler -
func SubscriptionL1Handler(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		L1Handler: L1Handler(msg.L1Handler),
	}
}

// L1Handler -
func L1Handler(model *storage.L1Handler) *pb.L1Handler {
	pbModel := &pb.L1Handler{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Status:     uint64(model.Status),
		Hash:       model.Hash,
		Position:   uint64(model.Position),
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		Selector:   model.Selector,
		Calldata:   model.CallData,
		MaxFee:     model.MaxFee.String(),
		Nonce:      model.Nonce.String(),
		Entrypoint: model.Entrypoint,
	}

	if model.ParsedCalldata != nil {
		parsed, err := json.Marshal(model.ParsedCalldata)
		if err != nil {
			return pbModel
		}
		pbModel.ParsedCalldata = parsed
	}
	return pbModel
}

// SubscriptionMessage -
func SubscriptionMessage(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Message: Message(msg.Message),
	}
}

// Message -
func Message(model *storage.Message) *pb.StarknetMessage {
	pbModel := &pb.StarknetMessage{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Contract:   model.Contract.Hash,
		From:       model.From.Hash,
		To:         model.To.Hash,
		Nonce:      model.Nonce.String(),
		Selector:   model.Selector,
		Payload:    model.Payload,
		ContractId: model.ContractID,
		FromId:     model.FromID,
		ToId:       model.ToID,
	}
	return pbModel
}

// SubscriptionStorageDiff -
func SubscriptionStorageDiff(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		StorageDiff: StorageDiff(msg.StorageDiff),
	}
}

// StorageDiff -
func StorageDiff(model *storage.StorageDiff) *pb.StorageDiff {
	pbModel := &pb.StorageDiff{
		Id:         model.ID,
		Height:     model.Height,
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		Key:        model.Key,
		Value:      model.Value,
	}
	return pbModel
}

// SubscriptionTokenBalance -
func SubscriptionTokenBalance(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		TokenBalance: TokenBalance(msg.TokenBalance),
	}
}

// TokenBalance -
func TokenBalance(model *storage.TokenBalance) *pb.TokenBalance {
	pbModel := &pb.TokenBalance{
		Owner:      model.Owner.Hash,
		OwnerId:    model.OwnerID,
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		TokenId:    model.TokenID.String(),
		Balance:    model.Balance.String(),
	}
	return pbModel
}

// SubscriptionTransfer -
func SubscriptionTransfer(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Transfer: Transfer(msg.Transfer),
	}
}

// Transfer -
func Transfer(model *storage.Transfer) *pb.Transfer {
	pbModel := &pb.Transfer{
		Id:         model.ID,
		Height:     model.Height,
		Time:       uint64(model.Time.Unix()),
		Contract:   model.Contract.Hash,
		ContractId: model.ContractID,
		From:       model.From.Hash,
		FromId:     model.FromID,
		To:         model.To.Hash,
		ToId:       model.ToID,
		Amount:     model.Amount.String(),
		TokenId:    model.TokenID.String(),
	}
	return pbModel
}

// SubscriptionToken -
func SubscriptionToken(id uint64, msg *subscriptions.Message) *pb.Subscription {
	return &pb.Subscription{
		Response: &generalPB.SubscribeResponse{
			Id: id,
		},
		Transfer: Transfer(msg.Transfer),
	}
}

// Token -
func Token(model *storage.Token) *pb.Token {
	pbModel := &pb.Token{
		Id:           model.ID,
		DeployHeight: model.DeployHeight,
		DeployTime:   uint64(model.DeployTime.Unix()),
		Contract:     model.Contract.Hash,
		ContractId:   model.ContractID,
		Owner:        model.Owner.Hash,
		Type:         int32(model.Type),
	}
	if model.Metadata != nil {
		metadata, err := json.Marshal(model.Metadata)
		if err != nil {
			return pbModel
		}
		pbModel.Metadata = metadata
	}
	return pbModel
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
