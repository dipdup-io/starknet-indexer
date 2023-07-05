package v0

import (
	"context"
	"errors"

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

type upgradeHandler func(data map[string]any, height uint64) ([]data.ProxyUpgrade, error)

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
		case "ModuleFunctionChange":
			handler = moduleFunctionChange
		case "implementation_upgraded":
			handler = implementationUpgraded
		default:
			continue
		}

		upgrades, err := handler(events[i].ParsedData, events[i].Height)
		if err != nil {
			return false, err
		}
		if len(upgrades) == 0 {
			return false, nil
		}
		if err := parser.resolver.UpgradeProxy(ctx, contract, upgrades, events[i].Height); err != nil {
			if errors.Is(err, resolver.ErrUnknownProxy) {
				return false, nil
			}
			return false, err
		}
		for i := range upgrades {
			key := data.NewProxyKey(contract.Hash, upgrades[i].Selector)
			txCtx.ProxyUpgrades.Add(key, struct{}{})
		}
		return true, nil
	}
	return false, nil
}

func (parser ProxyUpgrader) parseParams(ctx context.Context, txCtx data.TxContext, contract storage.Address, entrypoint string, params map[string]any) error {
	if len(params) == 0 {
		return nil
	}

	var handler upgradeHandler
	switch entrypoint {
	case "upgrade":
		handler = upgraded
	default:
		return nil
	}

	upgrades, err := handler(params, txCtx.Height)
	if err != nil {
		return err
	}
	if len(upgrades) == 0 {
		return nil
	}
	if err := parser.resolver.UpgradeProxy(ctx, contract, upgrades, txCtx.Height); err != nil {
		if errors.Is(err, resolver.ErrUnknownProxy) {
			return nil
		}
		return err
	}
	for i := range upgrades {
		key := data.NewProxyKey(contract.Hash, upgrades[i].Selector)
		txCtx.ProxyUpgrades.Add(key, struct{}{})
	}
	return err
}
