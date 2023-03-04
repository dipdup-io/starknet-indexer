package storage

import "github.com/dipdup-io/starknet-go-api/pkg/data"

// EntrypointType -
type EntrypointType int

const (
	EntrypointTypeUnknown EntrypointType = iota + 1
	EntrypointTypeExternal
	EntrypointTypeConstructor
)

// NewEntrypointType -
func NewEntrypointType(value string) EntrypointType {
	switch value {
	case data.EntrypointTypeExternal:
		return EntrypointTypeExternal
	case data.EntrypointTypeConstructor:
		return EntrypointTypeConstructor
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
	default:
		return Unknown
	}
}
