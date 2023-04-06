package main

import (
	indexerConfig "github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
)

// Config -
type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string               `yaml:"log_level" validate:"omitempty,oneof=debug trace info warn error fatal panic"`
	Indexer       indexerConfig.Config `yaml:"indexer"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
