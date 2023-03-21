package v0

import (
	"context"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/pkg/errors"
)

const (
	actualFeeContractHash = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
	actualFeeToHash       = "0x05dcd266a80b8a5f29f04d779c6b166b80150c24f2180a75e82427242dab20a9"
)

// FeeParser -
type FeeParser struct {
	cache    *cache.Cache
	resolver resolver.Resolver

	actualFeeContractId uint64
	actualFeeToId       uint64
}

// NewFeeParser -
func NewFeeParser(cache *cache.Cache, resolver resolver.Resolver) FeeParser {
	return FeeParser{
		cache:    cache,
		resolver: resolver,
	}
}

// ParseInvocation -
func (parser FeeParser) ParseActualFee(ctx context.Context, txCtx data.TxContext, actualFee starknetData.Felt) (*storage.Fee, error) {
	fee := actualFee.Decimal()
	if fee.IsZero() {
		return nil, nil
	}

	if parser.actualFeeContractId == 0 {
		address := storage.Address{
			Hash: starknetData.Felt(actualFeeContractHash).Bytes(),
		}
		if err := parser.resolver.FindAddress(ctx, &address); err != nil {
			return nil, err
		}
		parser.actualFeeContractId = address.ID
	}

	if parser.actualFeeToId == 0 {
		address := storage.Address{
			Hash: starknetData.Felt(actualFeeToHash).Bytes(),
		}
		if err := parser.resolver.FindAddress(ctx, &address); err != nil {
			return nil, err
		}
		parser.actualFeeToId = address.ID
	}

	return &storage.Fee{
		ID:         parser.resolver.NextTxId(),
		Height:     txCtx.Height,
		Time:       txCtx.Time,
		Status:     txCtx.Status,
		Amount:     fee,
		FromID:     txCtx.ContractId,
		ToID:       parser.actualFeeToId,
		ContractID: parser.actualFeeContractId,
	}, nil
}

// ParseInvocation -
func (parser FeeParser) ParseInvocation(ctx context.Context, txCtx data.TxContext, feeInvocation sequencer.Invocation) (*storage.Fee, error) {
	if len(feeInvocation.Events) == 0 && len(feeInvocation.InternalCalls) == 0 {
		return nil, nil
	}

	fee, err := parser.parseEvents(ctx, txCtx, feeInvocation.ContractAddress.Bytes(), feeInvocation.ClassHash.Bytes(), feeInvocation.Events)
	if err != nil {
		return nil, err
	}
	if fee != nil {
		return fee, nil
	}

	for i := range feeInvocation.InternalCalls {
		fee, err = parser.ParseInvocation(ctx, txCtx, feeInvocation.InternalCalls[i])
		if err != nil {
			return nil, err
		}
		if fee != nil {
			return fee, nil
		}
	}

	return nil, errors.New("can't parse fee transfer")
}

func (parser FeeParser) parseEvents(ctx context.Context, txCtx data.TxContext, contractHash, classHash []byte, events []starknetData.Event) (*storage.Fee, error) {
	abi, err := parser.cache.GetAbiByClassHash(ctx, classHash)
	if err != nil {
		return nil, err
	}

	contract, err := parser.cache.GetAddress(ctx, contractHash)
	if err != nil {
		return nil, err
	}

	for i := range events {
		parsed, name, err := decode.Event(parser.cache, abi, events[i].Keys, events[i].Data)
		if err != nil {
			return nil, err
		}

		if name != "Transfer" {
			continue
		}

		transfers, err := transfer(ctx, parser.resolver, txCtx, contract.ID, storage.Event{
			Height:     txCtx.Height,
			ParsedData: parsed,
		})
		if err != nil {
			return nil, err
		}
		if len(transfers) == 0 {
			return nil, err
		}

		return &storage.Fee{
			ID:         parser.resolver.NextTxId(),
			Height:     txCtx.Height,
			Time:       txCtx.Time,
			Status:     txCtx.Status,
			Amount:     transfers[0].Amount,
			FromID:     transfers[0].FromID,
			ToID:       transfers[0].ToID,
			ContractID: contract.ID,
		}, nil
	}

	return nil, nil
}
