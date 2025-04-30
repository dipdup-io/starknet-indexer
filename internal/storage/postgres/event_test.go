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

// EventTestSuite -
type EventTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
	pm            database.RangePartitionManager
}

// SetupSuite -
func (s *EventTestSuite) SetupSuite() {
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
	currentTime, err := time.Parse(time.RFC3339, "2021-11-20T13:00:35+00:00")
	s.Require().NoError(err)
	err = s.pm.CreatePartition(ctx, currentTime, storage.Event{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"fixtures/event.yml",
		),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *EventTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

// TestCountByContractID -
func (s *EventTestSuite) TestCountByContractID() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	eventCount, err := s.storage.Event.Count(ctx,
		[]storage.EventFilter{
			{
				Contract: storage.IdFilter{
					Eq: 14,
				},
			},
		},
	)
	s.Require().NoError(err)
	s.Require().EqualValues(int(eventCount), 2)
}

func TestSuiteEvent_Run(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
