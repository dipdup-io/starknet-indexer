package postgres

import (
	"context"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
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
	ERC20         models.IERC20
	ERC721        models.IERC721
	ERC1155       models.IERC1155
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
		ERC20:         NewERC20(strg.Connection()),
		ERC721:        NewERC721(strg.Connection()),
		ERC1155:       NewERC1155(strg.Connection()),
		TokenBalance:  NewTokenBalance(strg.Connection()),
		State:         NewState(strg.Connection()),

		PartitionManager: NewPartitionManager(strg.Connection()),
	}

	s.RollbackManager = NewRollbackManager(s.Transactable, s.State, s.Blocks)

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.PgGo) error {
	for _, data := range []storage.Model{
		&models.State{},
		&models.Address{},
		&models.Class{},
		&models.StorageDiff{},
		&models.Block{},
		&models.Invoke{},
		&models.Declare{},
		&models.Deploy{},
		&models.DeployAccount{},
		&models.L1Handler{},
		&models.Internal{},
		&models.Event{},
		&models.Message{},
		&models.Transfer{},
		&models.Fee{},
		&models.ERC20{},
		&models.ERC721{},
		&models.ERC1155{},
		&models.TokenBalance{},
		&models.Proxy{},
	} {
		if err := conn.DB().WithContext(ctx).Model(data).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		}); err != nil {
			if err := conn.Close(); err != nil {
				return err
			}
			return err
		}
	}
	return createIndices(ctx, conn)
}

func createIndices(ctx context.Context, conn *database.PgGo) error {
	return conn.DB().RunInTransaction(ctx, func(tx *pg.Tx) error {
		// Address
		if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS address_hash_idx ON address (hash)`); err != nil {
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
		return nil
	})
}
