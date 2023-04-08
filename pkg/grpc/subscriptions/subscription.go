package subscriptions

import (
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
}

// NewSubscription -
func NewSubscription(req *pb.SubscribeRequest) *Subscription {
	all := &Subscription{
		data: make(chan *pb.Subscription, 1024),
	}
	if req == nil {
		return all
	}
	all.blocks = req.Head

	if req.Declares != nil {
		all.declares = filters.NewDeclare(req.Declares)
	}
	if req.Deploys != nil {
		all.deploys = filters.NewDeploy(req.Deploys)
	}
	if req.DeployAccounts != nil {
		all.deployAccounts = filters.NewDeployAccount(req.DeployAccounts)
	}
	if req.Events != nil {
		all.events = filters.NewEvent(req.Events)
	}
	if req.Fees != nil {
		all.fees = filters.NewFee(req.Fees)
	}
	if req.Internals != nil {
		all.internals = filters.NewInternal(req.Internals)
	}
	if req.Invokes != nil {
		all.invokes = filters.NewInvoke(req.Invokes)
	}
	if req.L1Handlers != nil {
		all.l1Handlers = filters.NewL1Handler(req.L1Handlers)
	}
	if req.Msgs != nil {
		all.messages = filters.NewMessage(req.Msgs)
	}
	if req.StorageDiffs != nil {
		all.storageDiffs = filters.NewStorageDiff(req.StorageDiffs)
	}
	if req.TokenBalances != nil {
		all.tokenBalances = filters.NewTokenBalance(req.TokenBalances)
	}
	if req.Transfers != nil {
		all.transfers = filters.NewTransfer(req.Transfers)
	}

	return all
}

// Filter -
func (s *Subscription) Filter(msg *Message) bool {
	if msg == nil {
		return false
	}
	if msg.EndOfBlock {
		return true
	}
	if s.blocks {
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
