package api

import "github.com/dipdup-io/starknet-go-api/pkg/data"

func (block *SqdBlockResponse) GetDeployContractAddress(txIndex uint) data.Felt {
	for _, trace := range block.Traces {
		if trace.TraceType != data.TransactionTypeDeploy || trace.TraceType != data.TransactionTypeDeployAccount ||
			trace.TransactionIndex != txIndex {
			continue
		}

		return data.Felt(trace.ContractAddress)
	}
	return ""
}
