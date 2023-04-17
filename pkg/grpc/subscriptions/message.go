package subscriptions

import "github.com/dipdup-io/starknet-indexer/internal/storage"

// Message -
type Message struct {
	Block         *storage.Block
	Declare       *storage.Declare
	Deploy        *storage.Deploy
	DeployAccount *storage.DeployAccount
	Event         *storage.Event
	Fee           *storage.Fee
	Internal      *storage.Internal
	Invoke        *storage.Invoke
	L1Handler     *storage.L1Handler
	Message       *storage.Message
	StorageDiff   *storage.StorageDiff
	TokenBalance  *storage.TokenBalance
	Transfer      *storage.Transfer

	EndOfBlock bool
}

// NewBlockMessage -
func NewBlockMessage(block *storage.Block) *Message {
	return &Message{
		Block: block,
	}
}

// NewDeclareMessage -
func NewDeclareMessage(model *storage.Declare) *Message {
	return &Message{
		Declare: model,
	}
}

// NewDeployMessage -
func NewDeployMessage(model *storage.Deploy) *Message {
	return &Message{
		Deploy: model,
	}
}

// NewDeployAccountMessage -
func NewDeployAccountMessage(model *storage.DeployAccount) *Message {
	return &Message{
		DeployAccount: model,
	}
}

// NewEventMessage -
func NewEventMessage(model *storage.Event) *Message {
	return &Message{
		Event: model,
	}
}

// NewFeeMessage -
func NewFeeMessage(model *storage.Fee) *Message {
	return &Message{
		Fee: model,
	}
}

// NewInternalMessage -
func NewInternalMessage(model *storage.Internal) *Message {
	return &Message{
		Internal: model,
	}
}

// NewInvokeMessage -
func NewInvokeMessage(model *storage.Invoke) *Message {
	return &Message{
		Invoke: model,
	}
}

// NewL1HandlerMessage -
func NewL1HandlerMessage(model *storage.L1Handler) *Message {
	return &Message{
		L1Handler: model,
	}
}

// NewStarknetMessage -
func NewStarknetMessage(model *storage.Message) *Message {
	return &Message{
		Message: model,
	}
}

// NewStorageDiffMessage -
func NewStorageDiffMessage(model *storage.StorageDiff) *Message {
	return &Message{
		StorageDiff: model,
	}
}

// NewTokenBalanceMessage -
func NewTokenBalanceMessage(model *storage.TokenBalance) *Message {
	return &Message{
		TokenBalance: model,
	}
}

// NewTransferMessage -
func NewTransferMessage(model *storage.Transfer) *Message {
	return &Message{
		Transfer: model,
	}
}

// NewEndMessage -
func NewEndMessage() *Message {
	return &Message{
		EndOfBlock: true,
	}
}
