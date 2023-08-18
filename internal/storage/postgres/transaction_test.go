package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

func ptr[T any](v T) *T {
	return &v
}

// TransactionTest -
type TransactionTest struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
}

// SetupSuite -
func (s *TransactionTest) SetupSuite() {
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
func (s *TransactionTest) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *TransactionTest) TestAddresses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	defer tx.Close(ctx)

	addresses := []*storage.Address{
		{
			ID:      1,
			Hash:    []byte{0, 2, 3, 45},
			ClassID: nil,
			Height:  100,
		},
		{
			ID:      2,
			Hash:    []byte{0, 2, 3, 12},
			ClassID: ptr(uint64(1)),
			Height:  100,
		},
	}

	err = tx.SaveAddresses(ctx, addresses...)
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	response, err := s.storage.Address.List(ctx, 2, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(response, 2)
	s.Require().EqualValues(1, response[0].ID)
	s.Require().EqualValues(100, response[0].Height)
	s.Require().Nil(response[0].ClassID)
	s.Require().Equal([]byte{0, 2, 3, 45}, response[0].Hash)
	s.Require().NotNil(response[1].ClassID)
}

func (s *TransactionTest) TestClasses() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	defer tx.Close(ctx)

	classes := []*storage.Class{
		{
			ID:     1,
			Hash:   []byte{0, 2, 3, 45},
			Height: 100,
		},
		{
			ID:     2,
			Hash:   []byte{0, 2, 3, 12},
			Height: 100,
			Cairo:  1,
		},
	}

	err = tx.SaveClasses(ctx, classes...)
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	response, err := s.storage.Class.List(ctx, 2, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(response, 2)

	s.Require().EqualValues(1, response[0].ID)
	s.Require().EqualValues(100, response[0].Height)
	s.Require().EqualValues(0, response[0].Cairo)
	s.Require().Equal([]byte{0, 2, 3, 45}, response[0].Hash)

	s.Require().EqualValues(2, response[1].ID)
	s.Require().EqualValues(100, response[1].Height)
	s.Require().EqualValues(1, response[1].Cairo)
	s.Require().Equal([]byte{0, 2, 3, 12}, response[1].Hash)
}

func (s *TransactionTest) TestTokens() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/token.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	defer tx.Close(ctx)

	tokens := []*storage.Token{
		{
			FirstHeight: 100,
			ContractId:  1,
			TokenId:     decimal.NewFromInt(1),
			Type:        storage.TokenTypeERC1155,
		},
		{
			FirstHeight: 100,
			ContractId:  2,
			TokenId:     decimal.NewFromInt(2),
			Type:        storage.TokenTypeERC721,
		},
	}

	err = tx.SaveTokens(ctx, tokens...)
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	response, err := s.storage.Token.List(ctx, 2, 0, sdk.SortOrderDesc)
	s.Require().NoError(err)
	s.Require().Len(response, 2)

	s.Require().EqualValues(10001, response[1].ID)
	s.Require().EqualValues(100, response[1].FirstHeight)
	s.Require().EqualValues(storage.TokenTypeERC1155, response[1].Type)
	s.Require().EqualValues(1, response[1].ContractId)
	s.Require().EqualValues(1, response[1].TokenId.IntPart())

	s.Require().EqualValues(10002, response[0].ID)
	s.Require().EqualValues(100, response[0].FirstHeight)
	s.Require().EqualValues(storage.TokenTypeERC721, response[0].Type)
	s.Require().EqualValues(2, response[0].ContractId)
	s.Require().EqualValues(2, response[0].TokenId.IntPart())
}

func (s *TransactionTest) TestProxy() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	proxy := storage.Proxy{
		ID:         1,
		ContractID: 1,
		Hash:       []byte{1},
		Selector:   nil,
		EntityType: storage.EntityTypeClass,
		EntityID:   2,
		EntityHash: []byte{2},
	}

	err = tx.SaveProxy(ctx, proxy)
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	err = tx.Close(ctx)
	s.Require().NoError(err)

	response, err := s.storage.Proxy.GetByID(ctx, 1)
	s.Require().NoError(err)
	s.Require().EqualValues(1, response.ID)
	s.Require().EqualValues(1, response.ContractID)
	s.Require().Equal([]byte{1}, response.Hash)
	s.Require().Nil(response.Selector)
	s.Require().EqualValues(storage.EntityTypeClass, response.EntityType)
	s.Require().EqualValues(2, response.EntityID)
	s.Require().Equal([]byte{2}, response.EntityHash)

	tx2, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)

	err = tx2.DeleteProxy(ctx, proxy)
	s.Require().NoError(err)

	err = tx2.Flush(ctx)
	s.Require().NoError(err)

	err = tx2.Close(ctx)
	s.Require().NoError(err)

	_, err = s.storage.Proxy.GetByID(ctx, 1)
	s.Require().Error(err)
}

func (s *TransactionTest) TestTokenBalanceUpdate() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/token_balance.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	defer tx.Close(ctx)

	balances := []*storage.TokenBalance{
		{
			Balance:    decimal.RequireFromString("-1000000000000"),
			ContractID: 2816755,
			OwnerID:    2281170,
			TokenID:    decimal.NewFromInt(0),
		},
	}

	err = tx.SaveTokenBalanceUpdates(ctx, balances...)
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	response, err := s.storage.TokenBalance.Balances(ctx, 2816755, 0, 1, 0)
	s.Require().NoError(err)
	s.Require().Len(response, 1)

	s.Require().EqualValues(2816755, response[0].ContractID)
	s.Require().EqualValues(0, response[0].TokenID.IntPart())
	s.Require().EqualValues(0, response[0].Balance.IntPart())
	s.Require().EqualValues(2281170, response[0].OwnerID)
}

func (s *TransactionTest) TestUpdateStatus() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	pm := database.NewPartitionManager(s.storage.Connection(), database.PartitionByYear)
	currentTime, err := time.Parse(time.RFC3339, "2023-08-14T13:00:35+00:00")
	s.Require().NoError(err)
	err = pm.CreatePartition(ctx, currentTime, storage.DeployAccount{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files("fixtures/deploy_account.yml"),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())

	tx, err := BeginTransaction(ctx, s.storage.Transactable)
	s.Require().NoError(err)
	defer tx.Close(ctx)

	err = tx.UpdateStatus(ctx, 154957, storage.StatusNotReceived, &storage.DeployAccount{})
	s.Require().NoError(err)

	err = tx.Flush(ctx)
	s.Require().NoError(err)

	response, err := s.storage.DeployAccount.List(ctx, 10, 0, sdk.SortOrderAsc)
	s.Require().NoError(err)
	s.Require().Len(response, 10)

	for i := range response {
		if response[i].Height == 154957 {
			s.Require().EqualValues(storage.StatusNotReceived, response[i].Status)
		}
	}
}

func TestSuiteTransaction_Run(t *testing.T) {
	suite.Run(t, new(TransactionTest))
}
