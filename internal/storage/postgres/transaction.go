package postgres

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// Transaction -
type Transaction struct {
	storage.Transaction
}

// BeginTransaction -
func BeginTransaction(ctx context.Context, tx storage.Transactable) (Transaction, error) {
	t, err := tx.BeginTransaction(ctx)
	return Transaction{t}, err
}

// SaveAddress -
func (t Transaction) SaveAddresses(ctx context.Context, addresses ...*models.Address) error {
	if len(addresses) == 0 {
		return nil
	}

	values := t.Tx().NewValues(&addresses)
	_, err := t.Tx().NewInsert().
		With("_data", values).
		Model((*models.Address)(nil)).
		On("CONFLICT (hash) DO UPDATE").
		TableExpr("_data").
		Set("class_id = excluded.class_id").
		Set("height = excluded.height").
		Exec(ctx)
	return err
}

// SaveClasses -
func (t Transaction) SaveClasses(ctx context.Context, classes ...*models.Class) error {
	if len(classes) == 0 {
		return nil
	}

	values := t.Tx().NewValues(&classes)
	_, err := t.Tx().NewInsert().
		With("_data", values).
		Model((*models.Class)(nil)).
		On("CONFLICT (id) DO UPDATE").
		TableExpr("_data").
		Set("abi = excluded.abi").
		Set("type = excluded.type").
		Exec(ctx)
	return err
}

// SaveTokens -
func (t Transaction) SaveTokens(ctx context.Context, tokens ...*models.Token) error {
	if len(tokens) == 0 {
		return nil
	}
	_, err := t.Tx().
		NewInsert().
		Model(&tokens).
		On("CONFLICT ON CONSTRAINT token_unique_id DO NOTHING").
		Exec(ctx)
	return err
}

// SaveTokenBalanceUpdates -
func (t Transaction) SaveTokenBalanceUpdates(ctx context.Context, updates ...*models.TokenBalance) error {
	if len(updates) == 0 {
		return nil
	}

	values := t.Tx().NewValues(&updates)
	_, err := t.Tx().NewInsert().
		With("_data", values).
		Model((*models.TokenBalance)(nil)).
		On("CONFLICT (owner_id, contract_id, token_id) DO UPDATE").
		TableExpr("_data").
		Set("balance = token_balance.balance + excluded.balance").
		Exec(ctx)
	return err
}

// SaveProxy -
func (t Transaction) SaveProxy(ctx context.Context, proxy models.Proxy) error {
	values := t.Tx().NewValues(&proxy).Column("contract_id", "hash", "selector", "entity_type", "entity_id", "entity_hash")
	_, err := t.Tx().NewInsert().Column("contract_id", "hash", "selector", "entity_type", "entity_id", "entity_hash").
		With("_data", values).
		Model((*models.Proxy)(nil)).
		On("CONFLICT (hash, selector) DO UPDATE").
		TableExpr("_data").
		Set("entity_type = excluded.entity_type").
		Set("entity_id = excluded.entity_id").
		Set("entity_hash = excluded.entity_hash").
		Exec(ctx)
	return err
}

// DeleteProxy -
func (t Transaction) DeleteProxy(ctx context.Context, proxy models.Proxy) error {
	query := t.Tx().NewDelete().Model(&proxy).
		Where("hash = ?", proxy.Hash)

	if proxy.Selector != nil {
		query.Where("selector = ?", proxy.Selector)
	} else {
		query.Where("selector IS NULL")
	}

	_, err := query.Exec(ctx)
	return err
}

// UpdateStatus -
func (t Transaction) UpdateStatus(ctx context.Context, height uint64, status models.Status, model any) error {
	_, err := t.Tx().NewUpdate().Model(model).
		Where("height = ?", height).
		Set("status = ?", status).
		Exec(ctx)
	return err
}

// SaveClassReplaces -
func (t Transaction) SaveClassReplaces(ctx context.Context, replaces ...*models.ClassReplace) error {
	if len(replaces) == 0 {
		return nil
	}

	_, err := t.Tx().NewInsert().
		Model(&replaces).
		Returning("id").
		Exec(ctx)
	return err
}
