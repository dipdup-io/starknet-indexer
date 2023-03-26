package v0

import (
	"context"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
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

type upgradeHandler func(data map[string]any) ([]byte, error)

// Parse -
func (parser ProxyUpgrader) Parse(ctx context.Context, txCtx data.TxContext, contract storage.Address, events []storage.Event, entrypoint string, data map[string]any) error {
	ok, err := parser.parseEvents(ctx, txCtx, contract, events)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return parser.parseParams(ctx, txCtx, contract, entrypoint, data)
}

func (parser ProxyUpgrader) parseEvents(ctx context.Context, txCtx data.TxContext, contract storage.Address, events []storage.Event) (bool, error) {
	for i := range events {
		var handler upgradeHandler
		switch events[i].Name {
		case "Upgraded":
			handler = upgraded
		case "account_upgraded":
			handler = accountUpgraded
		default:
			continue
		}

		newImpl, err := handler(events[i].ParsedData)
		if err != nil {
			return false, err
		}
		if err := parser.resolver.UpgradeProxy(ctx, contract, newImpl, events[i].Height); err != nil {
			return false, err
		}
		contractAddress := encoding.EncodeHex(contract.Hash)
		txCtx.ProxyUpgrades[contractAddress] = struct{}{}
		return true, nil
	}
	return false, nil
}

func (parser ProxyUpgrader) parseParams(ctx context.Context, txCtx data.TxContext, contract storage.Address, entrypoint string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}

	var handler upgradeHandler
	switch entrypoint {
	case "upgrade":
		handler = upgraded
	default:
		return nil
	}

	newImpl, err := handler(data)
	if err != nil {
		return err
	}
	if newImpl == nil {
		return nil
	}
	if err := parser.resolver.UpgradeProxy(ctx, contract, newImpl, txCtx.Height); err != nil {
		return err
	}
	contractAddress := encoding.EncodeHex(contract.Hash)
	txCtx.ProxyUpgrades[contractAddress] = struct{}{}
	return err
}
