package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dipdup-io/starknet-indexer/pkg/types"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

// TransferTestSuite -
type TransferTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
	pm            database.RangePartitionManager
}

// SetupSuite -
func (s *TransferTestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := database.NewPostgreSQLContainer(ctx, database.PostgreSQLContainerConfig{
		User:     "user",
		Password: "password",
		Database: "db_test",
		Port:     5432,
		Image:    "postgres:15",
	})
	s.Require().NoError(err)
	s.psqlContainer = psqlContainer

	store, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	}, true)
	s.Require().NoError(err)
	s.storage = store

	s.pm = database.NewPartitionManager(s.storage.Connection(), database.PartitionByYear)
	currentTime, err := time.Parse(time.RFC3339, "2021-11-16T10:00:35+00:00")
	s.Require().NoError(err)
	err = s.pm.CreatePartition(ctx, currentTime, storage.Transfer{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"fixtures/transfer.yml",
			"fixtures/address.yml",
		),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *TransferTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TransferTestSuite) TestFilterByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	transfers, err := s.storage.Transfer.Filter(ctx, []storage.TransferFilter{
		{
			Height: storage.IntegerFilter{
				Eq: 1,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(transfers, 1)
}

func (s *TransferTestSuite) TestCountByFromAddress() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	fromAddress, err := types.HexFromString("0x000000000000000000000000c84dd7fd43a7defb5b7a15c4fbbe11cbba6db1ba")
	s.Require().NoError(err)

	transfersCount, err := s.storage.Transfer.Count(ctx, []storage.TransferFilter{
		{
			From: storage.BytesFilter{
				Eq: fromAddress,
			},
		},
	})
	s.Require().NoError(err)
	s.Require().EqualValues(int(transfersCount), 1)
}

func TestSuiteTransfer_Run(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
