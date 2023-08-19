package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

// StorageDiffTest -
type StorageDiffTest struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *StorageDiffTest) SetupSuite() {
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

	storage, err := Create(ctx, config.Database{
		Kind:     config.DBKindPostgres,
		User:     s.psqlContainer.Config.User,
		Database: s.psqlContainer.Config.Database,
		Password: s.psqlContainer.Config.Password,
		Host:     s.psqlContainer.Config.Host,
		Port:     s.psqlContainer.MappedPort().Int(),
	})
	s.Require().NoError(err)
	s.storage = storage
}

// TearDownSuite -
func (s *StorageDiffTest) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *StorageDiffTest) TestGetOnBlock() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/storage_diff.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	b, err := hex.DecodeString("01e2cd4b3588e8f6f9c4e89fb0e293bf92018c96d7a93ee367d29a284223b6ff")
	s.Require().NoError(err)

	sd, err := s.storage.StorageDiff.GetOnBlock(ctx, 0, 8, b)
	s.Require().NoError(err)
	s.Require().EqualValues(8, sd.ID)

	_, err = s.storage.StorageDiff.GetOnBlock(ctx, 0, 6, b)
	s.Require().Error(err)
}

func TestSuiteStorageDiff_Run(t *testing.T) {
	suite.Run(t, new(StorageDiffTest))
}
