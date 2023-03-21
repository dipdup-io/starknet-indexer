package helpers

import (
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
)

// NeedDecode -
func NeedDecode(calldata []string, invocation *sequencer.Invocation) bool {
	return len(calldata) > 0 || (invocation != nil && len(invocation.Events) > 0)
}
