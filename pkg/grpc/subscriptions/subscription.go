package subscriptions

import (
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Subscription -
type Subscription struct {
	data   chan *pb.Subscription
	blocks bool
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
	if msg.Block != nil && s.blocks {
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
