package adapter

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
	"sort"
)

func ConvertTraces(block *api.SqdBlockResponse) ([]starknet.Trace, error) {
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

		invocation := buildInvocationTree(txTraces)
		resultTrace.FunctionInvocation = &invocation
		resultTraces = append(resultTraces, resultTrace)
	}

	return resultTraces, nil
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

func buildInvocationTree(flatInvocations []api.TraceResponse) starknet.Invocation {
	var root starknet.Invocation
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
			InternalCalls:      make([]starknet.Invocation, 0),
			Events:             make([]data.Event, 0),
			Messages:           make([]data.Message, 0),
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
