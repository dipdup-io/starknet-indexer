package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

// DeclareTestSuite -
type DeclareTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
	pm            database.RangePartitionManager
}

// SetupSuite -
func (s *DeclareTestSuite) SetupSuite() {
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
	err = s.pm.CreatePartition(ctx, currentTime, storage.Declare{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/declare.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *DeclareTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *DeclareTestSuite) TestFilterByVersion() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	declares, err := s.storage.Declare.Filter(ctx, []storage.DeclareFilter{
		{
			Version: storage.EnumFilter{
				Eq: 2,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(declares, 3)
}

func (s *DeclareTestSuite) TestFilterByStatus() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	declares, err := s.storage.Declare.Filter(ctx, []storage.DeclareFilter{
		{
			Status: storage.EnumFilter{
				Eq: 6,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(declares, 3)
}

func (s *DeclareTestSuite) TestFilterByHeight() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	declares, err := s.storage.Declare.Filter(ctx, []storage.DeclareFilter{
		{
			Height: storage.IntegerFilter{
				Eq: 154930,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(declares, 1)
}

func (s *DeclareTestSuite) TestFilterById() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	declares, err := s.storage.Declare.Filter(ctx, []storage.DeclareFilter{
		{
			ID: storage.IntegerFilter{
				Gt: 394709233,
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(declares, 2)
}

func (s *DeclareTestSuite) TestFilterByTime() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	declares, err := s.storage.Declare.Filter(ctx, []storage.DeclareFilter{
		{
			Time: storage.TimeFilter{
				Between: &storage.BetweenFilter{
					From: 1692017700,
					To:   1692018000,
				},
			},
		},
	}, storage.WithLimitFilter(3))
	s.Require().NoError(err)
	s.Require().Len(declares, 2)
}

func TestSuiteDeclare_Run(t *testing.T) {
	suite.Run(t, new(DeclareTestSuite))
}
