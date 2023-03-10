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

	InvokeV0      models.IInvokeV0
	InvokeV1      models.IInvokeV1
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
		InvokeV0:      NewInvokeV0(strg.Connection()),
		InvokeV1:      NewInvokeV1(strg.Connection()),
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
	}

	return s, nil
}

func initDatabase(ctx context.Context, conn *database.PgGo) error {
	for _, data := range []storage.Model{
		&models.Address{},
		&models.Class{},
		&models.StorageDiff{},
		&models.Block{},
		&models.InvokeV0{},
		&models.InvokeV1{},
		&models.Declare{},
		&models.Deploy{},
		&models.DeployAccount{},
		&models.L1Handler{},
		&models.Internal{},
		&models.Event{},
		&models.Message{},
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
		return nil
	})
}
