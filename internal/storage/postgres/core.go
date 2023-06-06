package postgres

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rs/zerolog/log"
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
	StorageDiff   models.IStorageDiff
	Proxy         models.IProxy
	Transfer      models.ITransfer
	Fee           models.IFee
	Token         models.IToken
	TokenBalance  models.ITokenBalance
	State         models.IState

	PartitionManager PartitionManager
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
		StorageDiff:   NewStorageDiff(strg.Connection()),
		Proxy:         NewProxy(strg.Connection()),
		Transfer:      NewTransfer(strg.Connection()),
		Fee:           NewFee(strg.Connection()),
		Token:         NewToken(strg.Connection()),
		TokenBalance:  NewTokenBalance(strg.Connection()),
		State:         NewState(strg.Connection()),

		PartitionManager: NewPartitionManager(strg.Connection()),
	}

	s.RollbackManager = NewRollbackManager(s.Transactable, s.State, s.Blocks)

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.PgGo) error {
	for _, data := range models.Models {
		if err := conn.DB().WithContext(ctx).Model(data).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			if err := conn.Close(); err != nil {
				return err
			}
			return err
		}
	}

	if err := makeComments(ctx, conn); err != nil {
		return err
	}

	return createIndices(ctx, conn)
}

func createIndices(ctx context.Context, conn *database.PgGo) error {
	log.Info().Msg("creating indexes...")
	return conn.DB().RunInTransaction(ctx, func(tx *pg.Tx) error {
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

		// Storage diff
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS storage_diff_key_idx ON storage_diff (key)`); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS storage_diff_contract_id_idx ON storage_diff (contract_id)`); err != nil {
			return err
		}

		// Invoke
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS invoke_height_idx ON invoke USING BRIN (height)`); err != nil {
			return err
		}

		// Declare
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS declare_height_idx ON declare USING BRIN (height)`); err != nil {
			return err
		}

		// Deploy
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_height_idx ON deploy USING BRIN (height)`); err != nil {
			return err
		}

		// DeployAccount
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS deploy_account_height_idx ON deploy_account USING BRIN (height)`); err != nil {
			return err
		}

		// L1 handler
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS l1_handler_height_idx ON l1_handler USING BRIN (height)`); err != nil {
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
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS token_owner_idx ON token (owner_id)`); err != nil {
			return err
		}

		return nil
	})
}
