package indexer

import (
	"math/big"

	starknet "github.com/dipdup-io/starknet-go-api/pkg/api"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func getInternalModels(rawBlock starknet.BlockWithTxs) (models.Block, error) {
	block := models.Block{
		Height:           rawBlock.BlockNumber,
		Time:             rawBlock.Timestamp,
		Hash:             rawBlock.BlockHash,
		ParentHash:       rawBlock.ParentHash,
		NewRoot:          rawBlock.NewRoot,
		SequencerAddress: rawBlock.SequencerAddress,
		Status:           models.NewStatus(rawBlock.Status),
		TxCount:          len(rawBlock.Transactions),

		InvokeV0:      make([]models.InvokeV0, 0),
		InvokeV1:      make([]models.InvokeV1, 0),
		Declare:       make([]models.Declare, 0),
		Deploy:        make([]models.Deploy, 0),
		DeployAccount: make([]models.DeployAccount, 0),
		L1Handler:     make([]models.L1Handler, 0),
	}

	for i := range rawBlock.Transactions {
		switch typed := rawBlock.Transactions[i].Body.(type) {
		case *starknet.InvokeV0:
			block.InvokeV0 = append(block.InvokeV0, getInvokeV0(typed, block, rawBlock.Transactions[i].TransactionHash))
		case *starknet.InvokeV1:
			block.InvokeV1 = append(block.InvokeV1, getInvokeV1(typed, block, rawBlock.Transactions[i].TransactionHash))
		case *starknet.Declare:
			block.Declare = append(block.Declare, getDeclare(typed, block, rawBlock.Transactions[i].TransactionHash))
		case *starknet.Deploy:
			block.Deploy = append(block.Deploy, getDeploy(typed, block, rawBlock.Transactions[i].TransactionHash))
		case *starknet.DeployAccount:
			block.DeployAccount = append(block.DeployAccount, getDeployAccount(typed, block, rawBlock.Transactions[i].TransactionHash))
		case *starknet.L1Handler:
			block.L1Handler = append(block.L1Handler, getL1Handler(typed, block, rawBlock.Transactions[i].TransactionHash))
		default:
			return block, errors.Errorf("unknown transaction type: %s", rawBlock.Transactions[i].Type)
		}
	}

	block.InvokeV0Count = len(block.InvokeV0)
	block.InvokeV1Count = len(block.InvokeV1)
	block.DeclareCount = len(block.Declare)
	block.DeployCount = len(block.Deploy)
	block.DeployAccountCount = len(block.DeployAccount)
	block.L1HandlerCount = len(block.L1Handler)

	return block, nil
}

func decimalFromHex(s string) decimal.Decimal {
	i, _ := new(big.Int).SetString(s, 0)
	return decimal.NewFromBigInt(i, 0)
}

func getInvokeV0(raw *starknet.InvokeV0, block models.Block, hash string) models.InvokeV0 {
	return models.InvokeV0{
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               hash,
		ContractAddress:    raw.ContractAddress,
		EntrypointSelector: raw.EntrypointSelector,
		Signature:          raw.Signature,
		CallData:           raw.Calldata,
		MaxFee:             decimalFromHex(raw.MaxFee),
		Nonce:              decimalFromHex(raw.Nonce),
	}
}

func getInvokeV1(raw *starknet.InvokeV1, block models.Block, hash string) models.InvokeV1 {
	return models.InvokeV1{
		Height:        block.Height,
		Time:          block.Time,
		Status:        block.Status,
		Hash:          hash,
		SenderAddress: raw.SenderAddress,
		Signature:     raw.Signature,
		CallData:      raw.Calldata,
		MaxFee:        decimalFromHex(raw.MaxFee),
		Nonce:         decimalFromHex(raw.Nonce),
	}
}

func getDeclare(raw *starknet.Declare, block models.Block, hash string) models.Declare {
	return models.Declare{
		Height:        block.Height,
		Time:          block.Time,
		Status:        block.Status,
		Hash:          hash,
		SenderAddress: raw.SenderAddress,
		ClassHash:     raw.ClassHash,
		Signature:     raw.Signature,
		MaxFee:        decimalFromHex(raw.MaxFee),
		Nonce:         decimalFromHex(raw.Nonce),
	}
}

func getDeploy(raw *starknet.Deploy, block models.Block, hash string) models.Deploy {
	return models.Deploy{
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                hash,
		ContractAddressSalt: raw.ContractAddressSalt,
		ConstructorCalldata: raw.ConstructorCalldata,
		ClassHash:           raw.ClassHash,
	}
}

func getDeployAccount(raw *starknet.DeployAccount, block models.Block, hash string) models.DeployAccount {
	return models.DeployAccount{
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                hash,
		ContractAddressSalt: raw.ContractAddressSalt,
		ConstructorCalldata: raw.ConstructorCalldata,
		ClassHash:           raw.ClassHash,
		MaxFee:              decimalFromHex(raw.MaxFee),
		Nonce:               decimalFromHex(raw.Nonce),
		Signature:           raw.Signature,
	}
}

func getL1Handler(raw *starknet.L1Handler, block models.Block, hash string) models.L1Handler {
	return models.L1Handler{
		Height:             block.Height,
		Time:               block.Time,
		Status:             block.Status,
		Hash:               hash,
		ContractAddress:    raw.ContractAddress,
		EntrypointSelector: raw.EntrypointSelector,
		CallData:           raw.Calldata,
		Nonce:              decimalFromHex(raw.Nonce),
	}
}
