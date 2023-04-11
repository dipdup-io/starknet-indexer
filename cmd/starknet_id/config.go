package main

import (
	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
)

// Config -
type Config struct {
	LogLevel string             `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	GRPC     *grpc.ClientConfig `yaml:"grpc" validate:"required"`
}

// Substitute -
func (c *Config) Substitute() error {
	return nil
}
