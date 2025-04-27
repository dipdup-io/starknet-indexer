package postgres

import (
	"context"
	"database/sql"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/types"
	"github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/go-lib/database"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// InternalTestSuite -
type InternalTestSuite struct {
	suite.Suite
	psqlContainer *database.PostgreSQLContainer
	storage       Storage
	internalPm    database.RangePartitionManager
	deployPm      database.RangePartitionManager
}

// SetupSuite -
func (s *InternalTestSuite) SetupSuite() {
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

	s.internalPm = database.NewPartitionManager(s.storage.Connection(), database.PartitionByYear)
	s.deployPm = database.NewPartitionManager(s.storage.Connection(), database.PartitionByYear)
	currentTime, err := time.Parse(time.RFC3339, "2021-11-20T13:00:35+00:00")
	s.Require().NoError(err)

	err = s.internalPm.CreatePartition(ctx, currentTime, storage.Internal{}.TableName())
	s.Require().NoError(err)

	err = s.deployPm.CreatePartition(ctx, currentTime, storage.Deploy{}.TableName())
	s.Require().NoError(err)

	db, err := sql.Open("postgres", s.psqlContainer.GetDSN())
	s.Require().NoError(err)

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"fixtures/internal_tx.yml",
			"fixtures/deploy.yml",
			"fixtures/address.yml",
		),
	)
	s.Require().NoError(err)
	s.Require().NoError(fixtures.Load())
	s.Require().NoError(db.Close())
}

// TearDownSuite -
func (s *InternalTestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.storage.Close())
	s.Require().NoError(s.psqlContainer.Terminate(ctx))
}

func (s *InternalTestSuite) TestGetDeployedContracts() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	// Test for deployer with ID 2
	deployerHex := "0x020cfa74ee3564b4cd5435cdace0f9c4d43b939620e4a0bb5076105df0a626c6"
	deployerBytes, err := types.HexFromString(deployerHex)
	s.Require().NoError(err)

	contracts, err := s.storage.Internal.GetDeployedContracts(ctx, deployerBytes)
	s.Require().NoError(err)
	s.Require().Equal(2, len(contracts), "Should find 2 contracts")

	// Check contract addresses for deployer ID 2
	expectedContractAddresses := map[string]bool{
		"0x06538fdd3aa353af8a87f5fe77d1f533ea82815076e30a86d65b72d3eb4f0b80": true,
		"0x0327d34747122d7a40f4670265b098757270a449ec80c4871450fffdab7c2fa8": true,
	}
	s.Require().True(expectedContractAddresses[contracts[0].ContractAddress.String()], "Contract address should match expected")
	s.Require().True(expectedContractAddresses[contracts[1].ContractAddress.String()], "Contract address should match expected")
	s.Require().NotEqual(contracts[0].ContractAddress.String(), contracts[1].ContractAddress.String(), "Contract addresses should be different")

	// Test for deployer with ID 4
	deployerHex = "0x031c887d82502ceb218c06ebb46198da3f7b92864a8223746bc836dda3e34b52"
	deployerBytes, err = types.HexFromString(deployerHex)
	s.Require().NoError(err)

	contracts, err = s.storage.Internal.GetDeployedContracts(ctx, deployerBytes)
	s.Require().NoError(err)
	s.Require().Equal(2, len(contracts), "Should find 2 contracts")

	// Check contract addresses for deployer ID 4
	expectedContractAddresses = map[string]bool{
		"0x031c9cdb9b00cb35cf31c05855c0ec3ecf6f7952a1ce6e3c53c3455fcd75a280": true,
		"0x06ee3440b08a9c805305449ec7f7003f27e9f7e287b83610952ec36bdc5a6bae": true,
	}
	s.Require().True(expectedContractAddresses[contracts[0].ContractAddress.String()], "Contract address should match expected")
	s.Require().True(expectedContractAddresses[contracts[1].ContractAddress.String()], "Contract address should match expected")
	s.Require().NotEqual(contracts[0].ContractAddress.String(), contracts[1].ContractAddress.String(), "Contract addresses should be different")

	// Test for deployer with ID 6
	deployerHex = "0x031c9cdb9b00cb35cf31c05855c0ec3ecf6f7952a1ce6e3c53c3455fcd75a280"
	deployerBytes, err = types.HexFromString(deployerHex)
	s.Require().NoError(err)

	contracts, err = s.storage.Internal.GetDeployedContracts(ctx, deployerBytes)
	s.Require().NoError(err)
	s.Require().Equal(1, len(contracts), "Should find 1 contract")

	// Check contract address for deployer ID 6
	expectedContractAddress := "0x0735596016a37ee972c42adef6a3cf628c19bb3794369c65d2c82ba034aecf2c"
	s.Require().Equal(expectedContractAddress, contracts[0].ContractAddress.String(), "Contract address should match expected")

	// Check transaction hashes match those in deploy.yml
	expectedTxHashes := map[string]bool{
		"0x03c1c2fdf26d6edf12e639df4d237e862b0bb9933108b5af7398dfcb074c3d30": true,
		"0x0643bea7f02451b24540d36b235af30c3633d31dfaef1e9bf2ea125668f4710c": true,
		"0x0524d94af8066fc4831a94bd1ee9ddbd611617066fefabf69be7ddb963f417d3": true,
		"0x03aeb00f1ebf5f11d94aa069a5196521bd360ffb98c02b0a8a96a0c0bd6b8c49": true,
		"0x0114a5b2feb2226f7f25d650b58201a2602db3f59406b12d739c17936603ef85": true,
	}

	for _, contract := range contracts {
		s.Require().True(expectedTxHashes[contract.TxHash.String()], "Transaction hash should match expected")
		s.Require().NotEmpty(contract.Time, "Deploy time should not be empty")
	}
}

func TestSuiteInternal_Run(t *testing.T) {
	suite.Run(t, new(InternalTestSuite))
}
