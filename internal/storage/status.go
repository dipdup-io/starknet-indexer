package storage

import "github.com/dipdup-io/starknet-go-api/pkg/data"

// Status -
type Status int

const (
	StatusUnknown Status = iota + 1
	StatusNotReceived
	StatusReceived
	StatusPending
	StatusRejected
	StatusAcceptedOnL2
	StatusAcceptedOnL1
)

// NewStatus -
func NewStatus(value string) Status {
	switch value {
	case data.StatusNotReceived:
		return StatusNotReceived
	case data.StatusReceived:
		return StatusReceived
	case data.StatusPending:
		return StatusPending
	case data.StatusRejected:
		return StatusRejected
	case data.StatusAcceptedOnL2:
		return StatusAcceptedOnL2
	case data.StatusAcceptedOnL1:
		return StatusAcceptedOnL1
	default:
		return StatusUnknown
	}
}

// String -
func (s Status) String() string {
	switch s {
	case StatusNotReceived:
		return data.StatusNotReceived
	case StatusReceived:
		return data.StatusReceived
	case StatusPending:
		return data.StatusPending
	case StatusRejected:
		return data.StatusRejected
	case StatusAcceptedOnL2:
		return data.StatusAcceptedOnL2
	case StatusAcceptedOnL1:
		return data.StatusAcceptedOnL1
	default:
		return Unknown
	}
}
