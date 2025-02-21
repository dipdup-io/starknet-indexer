package adapter

import (
	"context"
	"fmt"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
	"time"
)

func (a *Adapter) convert(_ context.Context, block *api.SqdBlockResponse) (receiver.Result, error) {
	result := receiver.NewResult()
	b := receiver.Block{
		Height:           block.Header.Number,
		Status:           storage.NewStatus(block.Header.Status),
		Hash:             data.Felt(block.Header.Hash).Bytes(),
		ParentHash:       data.Felt(block.Header.ParentHash).Bytes(),
		NewRoot:          encoding.MustDecodeHex(block.Header.NewRoot),
		Time:             time.Unix(block.Header.Timestamp, 0).UTC(),
		SequencerAddress: encoding.MustDecodeHex(block.Header.SequencerAddress),
		Transactions:     ConvertTransactions(block),
		Receipts:         nil,
	}
	result.SetBlock(b)

	traces, err := ConvertTraces(block)
	if err != nil {
		return result, err
	}
	result.SetTraces(traces)

	stateUpdates, err := ConvertStateUpdates(block)
	if err != nil {
		return result, err
	}
	result.SetStateUpdates(stateUpdates)

	return result, nil
}

func uint64ToFelt(value *uint64) data.Felt {
	if value == nil {
		return ""
	}
	return data.Felt(fmt.Sprintf("%d", *value))
}

func stringToFelt(value *string) data.Felt {
	if value == nil {
		return ""
	}
	return data.Felt(*value)
}

func parseStringSlice(value *[]string) []string {
	if value == nil {
		return []string{}
	}
	return *value
}

func parseString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringSliceToFeltSlice(income []string) []data.Felt {
	result := make([]data.Felt, len(income))
	for i := range income {
		result[i] = data.Felt(income[i])
	}
	return result
}
