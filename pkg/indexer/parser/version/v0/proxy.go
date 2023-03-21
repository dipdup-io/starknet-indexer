package v0

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
)

// ProxyUpgrader -
type ProxyUpgrader struct {
	resolver resolver.Resolver
}

// NewProxyUpgrader -
func NewProxyUpgrader(resolver resolver.Resolver) ProxyUpgrader {
	return ProxyUpgrader{
		resolver: resolver,
	}
}

// Parse -
func (parser ProxyUpgrader) Parse(ctx context.Context, contract storage.Address, events []storage.Event) error {
	for i := range events {
		switch events[i].Name {
		case "Upgraded":
			if err := upgraded(ctx, parser.resolver, contract, events[i]); err != nil {
				return err
			}
		default:
		}
	}
	return nil
}
