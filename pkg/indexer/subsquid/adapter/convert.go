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

func (a *Adapter) convert(_ context.Context, block *api.SqdBlockResponse) error {
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
	result.Block = b

	ConvertTraces(block)
	ConvertStateUpdates(block)

	return nil
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
