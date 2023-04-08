package filters

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

// Message -
type Message struct {
	*pb.MessageFilter
	isEmpty bool
}

// NewMessage -
func NewMessage(req *pb.MessageFilter) Message {
	msg := Message{
		isEmpty: true,
	}
	if req == nil {
		return msg
	}
	msg.isEmpty = false
	msg.MessageFilter = req
	return msg
}

// Filter -
func (f Message) Filter(data storage.Message) bool {
	if f.isEmpty {
		return true
	}

	if !validInteger(f.Height, data.Height) {
		return false
	}

	if !validTime(f.Time, data.Time) {
		return false
	}

	// TODO: think about passing contract address
	if !validBytes(f.Contract, data.Contract.Hash) {
		return false
	}

	if !validBytes(f.From, data.From.Hash) {
		return false
	}

	if !validBytes(f.To, data.To.Hash) {
		return false
	}

	if !validEquality(f.Selector, data.Selector) {
		return false
	}

	return true
}
