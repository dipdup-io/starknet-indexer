package helpers

import (
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// NeedDecode -
func NeedDecode(status storage.Status, calldata []string, invocation *sequencer.Invocation) bool {
	return status != storage.StatusReverted &&
		(len(calldata) > 0 || (invocation != nil && len(invocation.Events) > 0))
}
