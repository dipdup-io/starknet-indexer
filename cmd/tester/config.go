package main

import (
	"github.com/dipdup-io/starknet-indexer/pkg/grpc"
	indexerConfig "github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
)

// Config -
type Config struct {
	config.Config `yaml:",inline"`
	GRPC          *grpc.ClientConfig   `yaml:"grpc" validate:"required"`
	LogLevel      string               `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	Indexer       indexerConfig.Config `yaml:"indexer"`
	GraphQlUrl    string               `yaml:"graphql_url" validate:"required,url"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
