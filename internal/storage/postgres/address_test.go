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

// AddressTestSuite -
type AddressTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *AddressTestSuite) SetupSuite() {
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
func (s *AddressTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *AddressTestSuite) TestGetByHash() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/address.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	b, err := hex.DecodeString("031c9cdb9b00cb35cf31c05855c0ec3ecf6f7952a1ce6e3c53c3455fcd75a280")
	s.Require().NoError(err)
	address, err := s.storage.Address.GetByHash(ctx, b)
	s.Require().NoError(err)
	s.Require().EqualValues(6, address.ID)
}

func (s *AddressTestSuite) TestGetAddresses() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/address.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Address.GetAddresses(ctx, 2, 4, 6, 7)
	s.Require().NoError(err)
	s.Require().Len(addresses, 3)
	s.Require().EqualValues(addresses[0].ID, 2)
	s.Require().EqualValues(addresses[1].ID, 4)
	s.Require().EqualValues(addresses[2].ID, 6)
}

func (s *AddressTestSuite) TestGetByHashes() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/address.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	b1, err := hex.DecodeString("0735596016a37ee972c42adef6a3cf628c19bb3794369c65d2c82ba034aecf2c")
	s.Require().NoError(err)
	b2, err := hex.DecodeString("031c9cdb9b00cb35cf31c05855c0ec3ecf6f7952a1ce6e3c53c3455fcd75a280")
	s.Require().NoError(err)
	b3, err := hex.DecodeString("000000000000000000000000c84dd7fd43a7defb5b7a15c4fbbe11cbba6db1ba")
	s.Require().NoError(err)

	hash := [][]byte{b1, b2, b3}

	ids, err := s.storage.Address.GetByHashes(ctx, hash)
	s.Require().NoError(err)
	s.Require().Len(ids, 3)
}

func (s *AddressTestSuite) TestFilter() {
	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/address.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	addresses, err := s.storage.Address.Filter(ctx, []storage.AddressFilter{
		{
			OnlyStarknet: true,
		}, {
			ID: storage.IntegerFilter{
				Gt: 11,
			},
		},
	}, storage.WithLimitFilter(10))
	s.Require().NoError(err)
	s.Require().Len(addresses, 9)
}

func TestSuiteAddress_Run(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}
