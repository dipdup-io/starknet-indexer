package adapter

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
)

func ConvertTransactions(block *api.SqdBlockResponse) []receiver.Transaction {
	txs := block.Transactions
	resultTxs := make([]receiver.Transaction, len(txs))
	for i, tx := range txs {
		var body any
		switch tx.Type {
		case data.TransactionTypeInvoke:
			body = data.Invoke{
				MaxFee:             stringToFelt(tx.MaxFee),
				Nonce:              uint64ToFelt(tx.Nonce),
				ContractAddress:    stringToFelt(tx.ContractAddress),
				EntrypointSelector: stringToFelt(tx.EntryPointSelector),
				SenderAddress:      stringToFelt(tx.SenderAddress),
				Signature:          parseStringSlice(tx.Signature),
				Calldata:           parseStringSlice(tx.Calldata),
			}
		case data.TransactionTypeDeclare:
			body = data.Declare{
				MaxFee:            stringToFelt(tx.MaxFee),
				Nonce:             uint64ToFelt(tx.Nonce),
				SenderAddress:     stringToFelt(tx.SenderAddress),
				ContractAddress:   stringToFelt(tx.ContractAddress),
				Signature:         parseStringSlice(tx.Signature),
				ClassHash:         stringToFelt(tx.ClassHash),
				CompiledClassHash: stringToFelt(tx.CompiledClassHash),
			}
		case data.TransactionTypeDeploy:
			body = data.Deploy{
				ContractAddressSalt: parseString(tx.ContractAddressSalt),
				ConstructorCalldata: parseStringSlice(tx.Calldata),
				ClassHash:           stringToFelt(tx.ClassHash),
				ContractAddress:     block.GetDeployContractAddress(tx.TransactionIndex),
			}
		case data.TransactionTypeDeployAccount:
			body = data.DeployAccount{
				MaxFee:              stringToFelt(tx.MaxFee),
				Nonce:               uint64ToFelt(tx.Nonce),
				ContractAddress:     stringToFelt(tx.ContractAddress),
				ContractAddressSalt: parseString(tx.ContractAddressSalt),
				ClassHash:           stringToFelt(tx.ClassHash),
				ConstructorCalldata: parseStringSlice(tx.ConstructorCalldata),
				Signature:           parseStringSlice(tx.Signature),
			}
		case data.TransactionTypeL1Handler:
			body = data.L1Handler{
				Nonce:              uint64ToFelt(tx.Nonce),
				ContractAddress:    stringToFelt(tx.ContractAddress),
				EntrypointSelector: stringToFelt(tx.EntryPointSelector),
				Calldata:           parseStringSlice(tx.Calldata),
			}
		default:
			return nil
		}

		resultTxs[i] = receiver.Transaction{
			Type:    tx.Type,
			Version: data.Felt(tx.Version),
			Hash:    data.Felt(tx.TransactionHash),
			Body:    body,
		}
	}

	return resultTxs
}
