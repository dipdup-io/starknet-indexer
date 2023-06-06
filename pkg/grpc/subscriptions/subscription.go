package subscriptions

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/subscriptions/filters"
)

// Subscription -
type Subscription struct {
	data chan *pb.Subscription

	blocks         bool
	declares       filters.Declare
	deploys        filters.Deploy
	deployAccounts filters.DeployAccount
	events         filters.Event
	fees           filters.Fee
	internals      filters.Internal
	invokes        filters.Invoke
	l1Handlers     filters.L1Handler
	messages       filters.Message
	storageDiffs   filters.StorageDiff
	tokenBalances  filters.TokenBalance
	transfers      filters.Transfer
	tokens         filters.Token
	addresses      filters.Address
}

// NewSubscription -
func NewSubscription(ctx context.Context, db postgres.Storage, req *pb.SubscribeRequest) (*Subscription, error) {
	all := &Subscription{
		data: make(chan *pb.Subscription, 1024),
	}
	if req == nil {
		return all, nil
	}
	all.blocks = req.Head

	if req.Addresses != nil {
		all.addresses = filters.NewAddress(req.Addresses)
	}
	if req.Declares != nil {
		all.declares = filters.NewDeclare(req.Declares)
	}
	if req.Deploys != nil {
		fltr, err := filters.NewDeploy(ctx, db.Class, req.Deploys)
		if err != nil {
			return nil, err
		}
		all.deploys = fltr
	}
	if req.DeployAccounts != nil {
		fltr, err := filters.NewDeployAccount(ctx, db.Class, req.DeployAccounts)
		if err != nil {
			return nil, err
		}
		all.deployAccounts = fltr
	}
	if req.Events != nil {
		fltr, err := filters.NewEvent(ctx, db.Address, req.Events)
		if err != nil {
			return nil, err
		}
		all.events = fltr
	}
	if req.Fees != nil {
		fltr, err := filters.NewFee(ctx, db.Address, db.Class, req.Fees)
		if err != nil {
			return nil, err
		}
		all.fees = fltr
	}
	if req.Internals != nil {
		fltr, err := filters.NewInternal(ctx, db.Address, db.Class, req.Internals)
		if err != nil {
			return nil, err
		}
		all.internals = fltr
	}
	if req.Invokes != nil {
		fltr, err := filters.NewInvoke(ctx, db.Address, req.Invokes)
		if err != nil {
			return nil, err
		}
		all.invokes = fltr
	}
	if req.L1Handlers != nil {
		fltr, err := filters.NewL1Handler(ctx, db.Address, req.L1Handlers)
		if err != nil {
			return nil, err
		}
		all.l1Handlers = fltr
	}
	if req.Msgs != nil {
		fltr, err := filters.NewMessage(ctx, db.Address, req.Msgs)
		if err != nil {
			return nil, err
		}
		all.messages = fltr
	}
	if req.StorageDiffs != nil {
		fltr, err := filters.NewStorageDiff(ctx, db.Address, req.StorageDiffs)
		if err != nil {
			return nil, err
		}
		all.storageDiffs = fltr
	}
	if req.TokenBalances != nil {
		fltr, err := filters.NewTokenBalance(ctx, db.Address, req.TokenBalances)
		if err != nil {
			return nil, err
		}
		all.tokenBalances = fltr
	}
	if req.Transfers != nil {
		fltr, err := filters.NewTransfer(ctx, db.Address, req.Transfers)
		if err != nil {
			return nil, err
		}
		all.transfers = fltr
	}
	if req.Tokens != nil {
		fltr, err := filters.NewToken(ctx, db.Address, req.Tokens)
		if err != nil {
			return nil, err
		}
		all.tokens = fltr
	}

	return all, nil
}

// Filter -
func (s *Subscription) Filter(msg *Message) bool {
	if msg == nil {
		return false
	}
	if msg.EndOfBlock != nil {
		return true
	}
	if msg.Address != nil && s.addresses.Filter(*msg.Address) {
		return true
	}
	if s.blocks && msg.Block != nil {
		return true
	}
	if msg.Declare != nil && s.declares.Filter(*msg.Declare) {
		return true
	}
	if msg.Deploy != nil && s.deploys.Filter(*msg.Deploy) {
		return true
	}
	if msg.DeployAccount != nil && s.deployAccounts.Filter(*msg.DeployAccount) {
		return true
	}
	if msg.Event != nil && s.events.Filter(*msg.Event) {
		return true
	}
	if msg.Fee != nil && s.fees.Filter(*msg.Fee) {
		return true
	}
	if msg.Internal != nil && s.internals.Filter(*msg.Internal) {
		return true
	}
	if msg.Invoke != nil && s.invokes.Filter(*msg.Invoke) {
		return true
	}
	if msg.L1Handler != nil && s.l1Handlers.Filter(*msg.L1Handler) {
		return true
	}
	if msg.Message != nil && s.messages.Filter(*msg.Message) {
		return true
	}
	if msg.StorageDiff != nil && s.storageDiffs.Filter(*msg.StorageDiff) {
		return true
	}
	if msg.TokenBalance != nil && s.tokenBalances.Filter(*msg.TokenBalance) {
		return true
	}
	if msg.Transfer != nil && s.transfers.Filter(*msg.Transfer) {
		return true
	}
	if msg.Token != nil && s.tokens.Filter(*msg.Token) {
		return true
	}

	return false
}

// Send -
func (s *Subscription) Send(msg *pb.Subscription) {
	s.data <- msg
}

// Close -
func (s *Subscription) Close() error {
	close(s.data)
	return nil
}

// Listen -
func (s *Subscription) Listen() <-chan *pb.Subscription {
	return s.data
}
