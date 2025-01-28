package adapter

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
	"sort"
)

func ConvertTraces(block *api.SqdBlockResponse) []starknet.Trace {
	resultTraces := make([]starknet.Trace, 0)

	for _, tx := range block.Transactions {
		resultTrace := starknet.Trace{
			RevertedError:         "",
			ValidateInvocation:    nil,
			FunctionInvocation:    nil,
			FeeTransferInvocation: nil,
			Signature:             nil,
			TransactionHash:       data.Felt(tx.TransactionHash),
		}
		txTraces := getTxTraces(block.Traces, tx.TransactionIndex)
		if len(txTraces) == 0 {
			resultTraces = append(resultTraces, resultTrace)
			continue
		}

		buildTree(txTraces)

		addressDepth := 1
		tracesByDepth := getTracesByDepth(txTraces, addressDepth)

		for traceIdx, traceInDepth := range tracesByDepth {
			calldata := make([]data.Felt, len(traceInDepth.Calldata))
			for i, cd := range traceInDepth.Calldata {
				calldata[i] = data.Felt(cd)
			}

			result := make([]data.Felt, len(traceInDepth.Result))
			for i, r := range traceInDepth.Result {
				result[i] = data.Felt(r)
			}

			invokation := starknet.Invocation{
				CallerAddress:      data.Felt(traceInDepth.CallerAddress),
				ContractAddress:    data.Felt(traceInDepth.ContractAddress),
				Calldata:           calldata,
				CallType:           parseString(traceInDepth.CallType),
				ClassHash:          stringToFelt(traceInDepth.ClassHash),
				Selector:           stringToFelt(traceInDepth.EntryPointSelector),
				EntrypointType:     parseString(traceInDepth.EntryPointType), // todo: wait for sqd re-index
				Result:             result,
				ExecutionResources: starknet.ExecutionResources{},
				InternalCalls:      nil,
				Events:             nil, // todo: events
				Messages:           nil, // todo: messages
			}

			if addressDepth == 1 {
				resultTrace.FeeTransferInvocation = &invokation
			} else {
				invokation.InternalCalls[traceIdx].InternalCalls[0].InternalCalls[0] = invokation
			}
			addressDepth += 1
		}
	}

	return resultTraces
}

func getTxTraces(traces []api.TraceResponse, txIndex uint) []api.TraceResponse {
	var result []api.TraceResponse
	for _, trace := range traces {
		if trace.TransactionIndex == txIndex {
			result = append(result, trace)
		}
	}
	return result
}

func getTxHashByIndex(txs []api.Transaction, txIndex uint) string {
	for _, tx := range txs {
		if tx.TransactionIndex == txIndex {
			return tx.TransactionHash
		}
	}
	return ""
}

func getTracesByDepth(traces []api.TraceResponse, depth int) []api.TraceResponse {
	var result []api.TraceResponse
	for _, trace := range traces {
		if len(trace.TraceAddress) == depth {
			result = append(result, trace)
		}
	}
	return result
}

func buildTree(flatInvocations []api.TraceResponse) starknet.Invocation {
	var root starknet.Invocation
	// TODO: don't sort?
	sort.Slice(flatInvocations, func(i, j int) bool {
		return compareTraceAddresses(flatInvocations[i].TraceAddress, flatInvocations[j].TraceAddress)
	})

	for _, inv := range flatInvocations {
		calldata := make([]data.Felt, len(inv.Calldata))
		for i, cd := range inv.Calldata {
			calldata[i] = data.Felt(cd)
		}

		result := make([]data.Felt, len(inv.Result))
		for i, r := range inv.Result {
			result[i] = data.Felt(r)
		}
		current := starknet.Invocation{
			CallerAddress:      data.Felt(inv.CallerAddress),
			ContractAddress:    data.Felt(inv.ContractAddress),
			Calldata:           calldata,
			CallType:           parseString(inv.CallType),
			ClassHash:          stringToFelt(inv.ClassHash),
			Selector:           stringToFelt(inv.EntryPointSelector),
			EntrypointType:     parseString(inv.EntryPointType),
			Result:             result,
			ExecutionResources: starknet.ExecutionResources{},
			InternalCalls:      nil,
			Events:             nil, // todo: events
			Messages:           nil, // todo: messages
		}

		level := len(inv.TraceAddress)
		if level == 1 {
			root = current
			continue
		}

		parentIndex := inv.TraceAddress[:level-1]
		parent := findParentByOrder(&root, parentIndex)
		if parent != nil {
			parent.InternalCalls = append(parent.InternalCalls, current)
		}
	}

	return root
}

func compareTraceAddresses(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func findParentByOrder(root *starknet.Invocation, traceAddress []int) *starknet.Invocation {
	current := root
	for i := 1; i < len(traceAddress); i++ {
		if current == nil || len(current.InternalCalls) == 0 {
			return nil
		}
		current = &current.InternalCalls[len(current.InternalCalls)-1]
	}
	return current
}
