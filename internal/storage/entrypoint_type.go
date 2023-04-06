package storage

import "github.com/dipdup-io/starknet-go-api/pkg/data"

// EntrypointType -
type EntrypointType int

const (
	EntrypointTypeUnknown EntrypointType = iota + 1
	EntrypointTypeExternal
	EntrypointTypeConstructor
	EntrypointTypeL1Handler
)

// NewEntrypointType -
func NewEntrypointType(value string) EntrypointType {
	switch value {
	case data.EntrypointTypeExternal:
		return EntrypointTypeExternal
	case data.EntrypointTypeConstructor:
		return EntrypointTypeConstructor
	case data.EntrypointTypeL1Handler:
		return EntrypointTypeL1Handler
	default:
		return EntrypointTypeUnknown
	}
}

// String -
func (s EntrypointType) String() string {
	switch s {
	case EntrypointTypeExternal:
		return data.EntrypointTypeExternal
	case EntrypointTypeConstructor:
		return data.EntrypointTypeConstructor
	case EntrypointTypeL1Handler:
		return data.EntrypointTypeL1Handler
	default:
		return Unknown
	}
}
