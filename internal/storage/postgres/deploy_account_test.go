package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

// DeployAccountTestSuite -
type DeployAccountTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
	pm            database.RangePartitionManager
}

// SetupSuite -
func (s *DeployAccountTestSuite) SetupSuite() {
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
	})
	s.Require().NoError(err)
	s.storage = store

	s.pm = database.NewPartitionManager(s.storage.Connection(), database.PartitionByYear)
	currentTime, err := time.Parse(time.RFC3339, "2023-08-14T13:00:35+00:00")
	s.Require().NoError(err)
	err = s.pm.CreatePartition(ctx, currentTime, storage.DeployAccount{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"fixtures/deploy_account.yml",
			"fixtures/class.yml",
		),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *DeployAccountTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *DeployAccountTestSuite) TestFilterByClass() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	b, err := hex.DecodeString("025ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918")
	s.Require().NoError(err)

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			Class: storage.BytesFilter{
				Eq: b,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(deploys, 3)
}

func (s *DeployAccountTestSuite) TestFilterByStatus() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			Status: storage.EnumFilter{
				Eq: 6,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(deploys, 3)
}

func (s *DeployAccountTestSuite) TestFilterByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			Height: storage.IntegerFilter{
				Eq: 154958,
			},
		},
	}, storage.WithLimitFilter(10))
	s.Require().NoError(err)
	s.Require().Len(deploys, 6)
}

func (s *DeployAccountTestSuite) TestFilterById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			ID: storage.IntegerFilter{
				Gt: 394803554,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(deploys, 1)
}

func (s *DeployAccountTestSuite) TestFilterByTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			Time: storage.TimeFilter{
				Between: &storage.BetweenFilter{
					From: 1692019800,
					To:   1692020100,
				},
			},
		},
	}, storage.WithLimitFilter(10))
	s.Require().NoError(err)
	s.Require().Len(deploys, 6)
}

func (s *DeployAccountTestSuite) TestFilterByParsedCalldata() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	deploys, err := s.storage.DeployAccount.Filter(ctx, []storage.DeployAccountFilter{
		{
			ParsedCalldata: map[string]string{
				"selector": "0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463",
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(deploys, 3)
}

func (s *DeployAccountTestSuite) TestHashByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	hash, err := s.storage.DeployAccount.HashByHeight(ctx, 154958)
	s.Require().NoError(err)
	s.Require().Len(hash, 32)
}

func TestSuiteDeployAccount_Run(t *testing.T) {
	suite.Run(t, new(DeployAccountTestSuite))
}
