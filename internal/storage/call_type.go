package storage

import "github.com/dipdup-io/starknet-go-api/pkg/data"

// unknown
const (
	Unknown = "UNKNOWN"
)

// CallType -
type CallType int

const (
	CallTypeUnknown CallType = iota + 1
	CallTypeCall
	CallTypeDelegate
)

// NewCallType -
func NewCallType(value string) CallType {
	switch value {
	case data.CallTypeCall:
		return CallTypeCall
	case data.CallTypeDelegate:
		return CallTypeDelegate
	default:
		return CallTypeUnknown
	}
}

// String -
func (s CallType) String() string {
	switch s {
	case CallTypeCall:
		return data.CallTypeCall
	case CallTypeDelegate:
		return data.CallTypeDelegate
	default:
		return Unknown
	}
}
