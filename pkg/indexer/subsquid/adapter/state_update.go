package adapter

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
)

func ConvertStateUpdates(block *api.SqdBlockResponse) (data.StateUpdate, error) {
	storageDiffs := make(map[data.Felt][]data.KeyValue)
	for i := range block.StorageDiffs {
		address := data.Felt(block.StorageDiffs[i].Address)
		storageDiffs[address] = append(storageDiffs[address], data.KeyValue{
			Key:   data.Felt(block.StorageDiffs[i].Key),
			Value: data.Felt(block.StorageDiffs[i].Value),
		})
	}

	declaredClasses := make([]data.DeclaredClass, 0)
	for i := range block.StateUpdates[0].DeclaredClasses {
		declaredClass := block.StateUpdates[0].DeclaredClasses[i]
		declaredClasses = append(declaredClasses, data.DeclaredClass{
			ClassHash:         data.Felt(declaredClass.ClassHash),
			CompiledClassHash: data.Felt(declaredClass.CompiledClassHash),
		})
	}

	replacedClasses := make([]data.ReplacedClass, 0)
	for i := range block.StateUpdates[0].ReplacedClasses {
		replacedClass := block.StateUpdates[0].ReplacedClasses[i]
		replacedClasses = append(replacedClasses, data.ReplacedClass{
			Address:   data.Felt(replacedClass.ContractAddress),
			ClassHash: data.Felt(replacedClass.ClassHash),
		})
	}

	oldDeclaredContracts := make([]data.Felt, 0)
	for _, oldDeclaredContract := range block.StateUpdates[0].DeprecatedDeclaredClasses {
		oldDeclaredContracts = append(oldDeclaredContracts, data.Felt(oldDeclaredContract))
	}

	deployedContracts := make([]data.DeployedContract, 0)
	for i := range block.StateUpdates[0].DeployedContracts {
		deployedContract := block.StateUpdates[0].DeployedContracts[i]
		deployedContracts = append(deployedContracts, data.DeployedContract{
			Address:   data.Felt(deployedContract.Address),
			ClassHash: data.Felt(deployedContract.ClassHash),
		})
	}

	nonces := make(map[data.Felt]data.Felt)
	for i := range block.StateUpdates[0].Nonces {
		nonce := block.StateUpdates[0].Nonces[i]
		nonces[data.Felt(nonce.ContractAddress)] = data.Felt(nonce.Nonce)
	}

	stateUpdate := data.StateUpdate{
		BlockHash: data.Felt(block.Header.Hash),
		NewRoot:   data.Felt(block.StateUpdates[0].NewRoot),
		OldRoot:   data.Felt(block.StateUpdates[0].OldRoot),
		StateDiff: data.StateDiff{
			StorageDiffs:         storageDiffs,
			DeclaredClasses:      declaredClasses,
			ReplacedClasses:      replacedClasses,
			OldDeclaredContracts: oldDeclaredContracts,
			DeployedContracts:    deployedContracts,
			Nonces:               nonces,
		},
	}

	return stateUpdate, nil
}
