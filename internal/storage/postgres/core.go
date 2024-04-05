package postgres

import (
	"context"
	"database/sql"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

// Storage -
type Storage struct {
	*postgres.Storage

	Blocks models.IBlock

	Invoke        models.IInvoke
	Declare       models.IDeclare
	Deploy        models.IDeploy
	DeployAccount models.IDeployAccount
	L1Handler     models.IL1Handler
	Internal      models.IInternal
	Message       models.IMessage
	Event         models.IEvent
	Address       models.IAddress
	Class         models.IClass
	ClassReplaces models.IClassReplace
	StorageDiff   models.IStorageDiff
	Proxy         models.IProxy
	ProxyUpgrade  models.IProxyUpgrade
	Transfer      models.ITransfer
	Fee           models.IFee
	Token         models.IToken
	TokenBalance  models.ITokenBalance
	State         models.IState

	PartitionManager database.RangePartitionManager
	RollbackManager  RollbackManager
}

// Create -
func Create(ctx context.Context, cfg config.Database) (Storage, error) {
	strg, err := postgres.Create(ctx, cfg, initDatabase)
	if err != nil {
		return Storage{}, err
	}

	s := Storage{
		Storage:       strg,
		Blocks:        NewBlocks(strg.Connection()),
		Invoke:        NewInvoke(strg.Connection()),
		Declare:       NewDeclare(strg.Connection()),
		Deploy:        NewDeploy(strg.Connection()),
		DeployAccount: NewDeployAccount(strg.Connection()),
		L1Handler:     NewL1Handler(strg.Connection()),
		Internal:      NewInternal(strg.Connection()),
		Message:       NewMessage(strg.Connection()),
		Event:         NewEvent(strg.Connection()),
		Address:       NewAddress(strg.Connection()),
		Class:         NewClass(strg.Connection()),
		ClassReplaces: NewClassReplace(strg.Connection()),
		StorageDiff:   NewStorageDiff(strg.Connection()),
		Proxy:         NewProxy(strg.Connection()),
		ProxyUpgrade:  NewProxyUpgrade(strg.Connection()),
		Transfer:      NewTransfer(strg.Connection()),
		Fee:           NewFee(strg.Connection()),
		Token:         NewToken(strg.Connection()),
		TokenBalance:  NewTokenBalance(strg.Connection()),
		State:         NewState(strg.Connection()),

		PartitionManager: database.NewPartitionManager(strg.Connection(), database.PartitionByMonth),
	}

	s.RollbackManager = NewRollbackManager(s.Transactable, s.State, s.Blocks, s.ProxyUpgrade, s.ClassReplaces, s.Transfer)

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.Bun) error {
	if err := createTypes(ctx, conn); err != nil {
		return errors.Wrap(err, "creating custom types")
	}

	if err := database.CreateTables(ctx, conn, models.ModelsAny...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return err
	}

	if err := database.MakeComments(ctx, conn, models.ModelsAny...); err != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		return errors.Wrap(err, "make comments")
	}

	return createIndices(ctx, conn)
}

func createIndices(ctx context.Context, conn *database.Bun) error {
	log.Info().Msg("creating indexes...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Address
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS address_hash_idx ON address (hash)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS address_height_idx ON address USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS address_class_id_idx ON address (class_id)`); err != nil {
			return err
		}

		// Proxy
		if _, err := tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS proxy_hash_selector_idx ON proxy (hash, selector) NULLS NOT DISTINCT`); err != nil {
			return err
		}

		// ProxyUpgrade
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS proxy_upgrade_hash_selector_idx ON proxy_upgrade (hash, selector) NULLS NOT DISTINCT`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS proxy_upgrade_height_idx ON proxy_upgrade USING BRIN (height)`); err != nil {
			return err
		}

		// Storage diff
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS storage_diff_key_idx ON storage_diff (key)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS storage_diff_contract_id_idx ON storage_diff (contract_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS storage_diff_height_idx ON storage_diff USING BRIN (height)`); err != nil {
			return err
		}

		// Invoke
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS invoke_height_idx ON invoke USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS invoke_contract_id_idx ON invoke (contract_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS invoke_status_idx ON invoke (status)`); err != nil {
			return err
		}

		// Declare
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS declare_height_idx ON declare USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS declare_status_idx ON declare (status)`); err != nil {
			return err
		}

		// Deploy
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_height_idx ON deploy USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_status_idx ON deploy (status)`); err != nil {
			return err
		}

		// DeployAccount
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_account_height_idx ON deploy_account USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_account_status_idx ON deploy_account (status)`); err != nil {
			return err
		}

		// L1 handler
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS l1_handler_height_idx ON l1_handler USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS l1_handler_status_idx ON l1_handler (status)`); err != nil {
			return err
		}

		// Fee
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS fee_height_idx ON fee USING BRIN (height)`); err != nil {
			return err
		}

		// Internal
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS internal_tx_height_idx ON internal_tx USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS internal_tx_hash_idx ON internal_tx (hash)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS internal_tx_status_idx ON internal_tx (status)`); err != nil {
			return err
		}

		// Event
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS event_height_idx ON event USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS event_contract_id_idx ON event (contract_id, id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS event_name_idx ON event (name, id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS event_contract_name_idx ON event (contract_id, name, id)`); err != nil {
			return err
		}

		// Message
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS message_height_idx ON message USING BRIN (height)`); err != nil {
			return err
		}

		// Token
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS token_contract_idx ON token (contract_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS token_identity_idx ON token (contract_id, token_id)`); err != nil {
			return err
		}

		// Transfer
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS transfer_height_idx ON transfer USING BRIN (height)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS transfer_from_id_idx ON transfer (from_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS transfer_to_id_idx ON transfer (to_id)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS transfer_contract_id_idx ON transfer (contract_id)`); err != nil {
			return err
		}

		// Class Replace
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS class_replace_height_idx ON class_replace USING BRIN (height)`); err != nil {
			return err
		}

		// Class
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS class_height_idx ON class USING BRIN (height)`); err != nil {
			return err
		}

		// Block
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS block_height_idx ON block USING BRIN (height)`); err != nil {
			return err
		}

		return nil
	})
}

func createTypes(ctx context.Context, conn *database.Bun) error {
	log.Info().Msg("creating custom types...")
	return conn.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.ExecContext(
			ctx,
			`DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'token_type') THEN
					CREATE TYPE token_type AS ENUM ('erc20', 'erc721', 'erc1155');
				END IF;
			END$$;`,
		); err != nil {
			return err
		}
		return nil
	})
}
