package storage

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
	case "NOT_RECEIVED":
		return StatusNotReceived
	case "RECEIVED":
		return StatusReceived
	case "PENDING":
		return StatusPending
	case "REJECTED":
		return StatusRejected
	case "ACCEPTED_ON_L2":
		return StatusAcceptedOnL2
	case "ACCEPTED_ON_L1":
		return StatusAcceptedOnL1
	default:
		return StatusUnknown
	}
}

// String -
func (s Status) String() string {
	switch s {
	case StatusNotReceived:
		return "NOT_RECEIVED"
	case StatusReceived:
		return "RECEIVED"
	case StatusPending:
		return "PENDING"
	case StatusRejected:
		return "REJECTED"
	case StatusAcceptedOnL2:
		return "ACCEPTED_ON_L2"
	case StatusAcceptedOnL1:
		return "ACCEPTED_ON_L1"
	default:
		return "UNKNOWN"
	}
}
