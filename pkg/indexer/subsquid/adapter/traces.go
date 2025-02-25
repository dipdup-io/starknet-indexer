package adapter

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
	"golang.org/x/exp/slices"
	"sort"
)

func ConvertTraces(block *api.SqdBlockResponse) ([]starknet.Trace, error) {
	traces := make([]starknet.Trace, 0)

	for i := range block.Transactions {
		tx := block.Transactions[i]
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
			traces = append(traces, resultTrace)
			continue
		}
		txEvents := getTxEvents(block.Events, tx.TransactionIndex)
		txMessages := getTxMessages(block.Messages, tx.TransactionIndex)

		trace := buildTraceTree(txTraces, txEvents, txMessages)
		trace.TransactionHash = data.Felt(tx.TransactionHash)
		traces = append(traces, trace)
	}

	return traces, nil
}

func getTxTraces(traces []api.TraceResponse, txIndex uint) []api.TraceResponse {
	var result []api.TraceResponse
	for i := range traces {
		trace := traces[i]
		if trace.TransactionIndex == txIndex {
			result = append(result, trace)
		}
	}
	return result
}

func getTxEvents(events []api.Event, txIndex uint) []api.Event {
	var result []api.Event
	for i := range events {
		event := events[i]
		if event.TransactionIndex == txIndex {
			result = append(result, event)
		}
	}
	return result
}

func getTxMessages(messages []api.Message, txIndex uint) []api.Message {
	var result []api.Message
	for i := range messages {
		message := messages[i]
		if message.TransactionIndex == txIndex {
			result = append(result, message)
		}
	}
	return result
}

func buildTraceTree(flatInvocations []api.TraceResponse, events []api.Event, messages []api.Message) starknet.Trace {
	resultTrace := starknet.Trace{
		RevertedError:         "",
		ValidateInvocation:    nil,
		FunctionInvocation:    nil,
		FeeTransferInvocation: nil,
		Signature:             nil,
		TransactionHash:       "",
	}
	mapAddressInvocationType := make(map[int]string)
	sort.Slice(flatInvocations, func(i, j int) bool {
		res := compareTraceAddresses(flatInvocations[i].TraceAddress, flatInvocations[j].TraceAddress)
		return res
	})

	for invokationIndex := range flatInvocations {
		invocation := flatInvocations[invokationIndex]
		calldata := stringSliceToFeltSlice(invocation.Calldata)
		result := stringSliceToFeltSlice(invocation.Result)
		sqdEvents := filterEventsByAddress(events, invocation.TraceAddress)
		adaptedEvents := make([]data.Event, len(sqdEvents))

		for i := range sqdEvents {
			event := sqdEvents[i]
			keys := stringSliceToFeltSlice(event.Keys)
			eventData := stringSliceToFeltSlice(event.Data)

			eventOrder := uint64(event.EvenIndex)
			switch invocation.InvocationType {
			case "execute", "constructor":
			default:
				eventOrder = 0
			}

			adaptedEvents[i] = data.Event{
				Order:       eventOrder,
				FromAddress: "",
				Keys:        keys,
				Data:        eventData,
			}
		}

		sqdMessages := filterMessagesByAddress(messages, invocation.TraceAddress)
		adaptedMessages := make([]data.Message, len(sqdMessages))
		for i := range sqdMessages {
			message := sqdMessages[i]
			payload := stringSliceToFeltSlice(message.Payload)
			adaptedMessages[i] = data.Message{
				Order:       uint64(message.Order),
				FromAddress: parseString(message.FromAddress),
				ToAddress:   message.ToAddress,
				Selector:    "",
				Payload:     payload,
				Nonce:       "",
			}
		}

		currentInvocation := starknet.Invocation{
			CallerAddress:      data.Felt(invocation.CallerAddress),
			ContractAddress:    data.Felt(invocation.ContractAddress),
			Calldata:           calldata,
			CallType:           parseString(invocation.CallType),
			ClassHash:          stringToFelt(invocation.ClassHash),
			Selector:           stringToFelt(invocation.EntryPointSelector),
			EntrypointType:     parseString(invocation.EntryPointType),
			Result:             result,
			ExecutionResources: starknet.ExecutionResources{},
			InternalCalls:      make([]starknet.Invocation, 0),
			Events:             adaptedEvents,
			Messages:           adaptedMessages,
		}

		level := len(invocation.TraceAddress)
		if level == 1 {
			mapAddressInvocationType[invocation.TraceAddress[0]] = invocation.InvocationType
			switch invocation.InvocationType {
			case "fee_transfer":
				resultTrace.FeeTransferInvocation = &currentInvocation
			case "validate":
				resultTrace.ValidateInvocation = &currentInvocation
			case "execute", "constructor":
				if invocation.RevertReason != nil {
					resultTrace.RevertedError = parseString(invocation.RevertReason)
				} else {
					resultTrace.FunctionInvocation = &currentInvocation
				}
			}

			continue
		}

		parentIndex := invocation.TraceAddress[:level-1]
		parent := findParentByOrder(&resultTrace, parentIndex, mapAddressInvocationType)
		if parent != nil {
			parent.InternalCalls = append(parent.InternalCalls, currentInvocation)
		}
	}

	return resultTrace
}

func compareTraceAddresses(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func findParentByOrder(trace *starknet.Trace, traceAddress []int, mapAddressInvocationType map[int]string) *starknet.Invocation {
	var rootIndex int
	if len(traceAddress) == 1 {
		rootIndex = traceAddress[0]
	} else {
		rootIndex = traceAddress[:1][0]
	}

	var rootInvocation *starknet.Invocation
	switch invocationType := mapAddressInvocationType[rootIndex]; invocationType {
	case "fee_transfer":
		rootInvocation = trace.FeeTransferInvocation
	case "validate":
		rootInvocation = trace.ValidateInvocation
	case "execute", "constructor":
		rootInvocation = trace.FunctionInvocation
	}

	current := rootInvocation
	for i := 1; i < len(traceAddress); i++ {
		if current == nil || len(current.InternalCalls) == 0 {
			return nil
		}
		current = &current.InternalCalls[len(current.InternalCalls)-1]
	}
	return current
}

func filterEventsByAddress(events []api.Event, targetAddress []int) []api.Event {
	var result []api.Event

	for _, event := range events {
		if slices.Equal(event.TraceAddress, targetAddress) {
			result = append(result, event)
		}
	}
	return result
}

func filterMessagesByAddress(messages []api.Message, targetAddress []int) []api.Message {
	var result []api.Message

	for _, message := range messages {
		if slices.Equal(message.TraceAddress, targetAddress) {
			result = append(result, message)
		}
	}
	return result
}
