package v0

import (
	"context"
	"errors"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/shopspring/decimal"
)

// TransferParser -
type TransferParser struct {
	resolver resolver.Resolver
}

// NewTransferParser -
func NewTransferParser(resolver resolver.Resolver) TransferParser {
	return TransferParser{resolver: resolver}
}

// ParseCalldata -
func (parser TransferParser) ParseCalldata(
	ctx context.Context,
	txCtx data.TxContext,
	entrypoint string,
	calldata map[string]any,
) ([]storage.Transfer, error) {

	switch entrypoint {
	case "transfer":
		return parser.parseTransferCalldata(ctx, txCtx, calldata)
	case "transferFrom":
		if _, ok := calldata["from_"]; ok {
			return parser.parseTransferFromERC721Calldata(ctx, txCtx, calldata)
		}
		return parser.parseTransferFromCalldata(ctx, txCtx, calldata)
	case "safeTransferFrom":
		return parser.parseTransferFromERC721Calldata(ctx, txCtx, calldata)
	case "mint":
		if _, ok := calldata["amount"]; ok {
			return parser.parseMintErc20(ctx, txCtx, calldata)
		}
	}

	return nil, nil
}

// ParseEvents -
func (parser TransferParser) ParseEvents(
	ctx context.Context,
	txCtx data.TxContext,
	contract storage.Address,
	events []storage.Event,
) ([]storage.Transfer, error) {
	contractId := contract.ID
	if txCtx.ProxyId > 0 {
		contractId = txCtx.ProxyId
	}

	var (
		transfers = make([]storage.Transfer, 0)
		err       error
	)

	for i := range events {
		if events[i].ParsedData == nil {
			continue
		}

		var t []storage.Transfer

		switch events[i].Name {

		case "Transfer":
			if _, ok := events[i].ParsedData["tokenId"]; ok {
				t, err = transferERC721(ctx, parser.resolver, txCtx, contractId, events[i])
			} else {
				t, err = transfer(ctx, parser.resolver, txCtx, contractId, events[i])
			}
		case "TransferSingle":
			t, err = transferSingle(ctx, parser.resolver, txCtx, contractId, events[i])
		case "TransferBatch":
			t, err = transferBatch(ctx, parser.resolver, txCtx, contractId, events[i])
		case "deposit_handled":
			if events[i].Order > 0 {
				// if deposit_handled wasn't first event in transaction than Transfer was first
				continue
			}
			t, err = depositHandled(ctx, parser.resolver, txCtx, contract, events[i])
		case "withdraw_initiated":
			if events[i].Order > 0 {
				// if withdraw_initiated wasn't first event in transaction than Transfer was first
				continue
			}
			t, err = withdrawInitiated(ctx, parser.resolver, txCtx, contract, events[i])
		default:
			continue
		}
		if err != nil {
			if errors.Is(err, errInvalidTransfer) {
				continue
			}
			return nil, err
		}
		if len(t) > 0 {
			transfers = append(transfers, t...)
		}
	}
	return transfers, nil
}

func (parser TransferParser) parseTransferCalldata(ctx context.Context, txCtx data.TxContext, calldata map[string]any) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		InternalID:      txCtx.InternalID,
		FeeID:           txCtx.FeeID,
		ContractID:      txCtx.ContractId,
	}

	switch {
	case txCtx.Internal != nil:
		transfer.FromID = txCtx.Internal.CallerID
	case txCtx.Fee != nil:
		transfer.FromID = txCtx.Fee.CallerID
	default:
		transfer.FromID = txCtx.ContractId
	}

	recipientId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, "recipient")
	if err != nil {
		return nil, err
	}
	transfer.ToID = recipientId

	amount, err := parseTransferDecimal(calldata, "amount")
	if err != nil {
		return nil, err
	}
	transfer.Amount = amount

	return []storage.Transfer{transfer}, nil
}

func (parser TransferParser) parseTransferFromCalldata(ctx context.Context, txCtx data.TxContext, calldata map[string]any) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		InternalID:      txCtx.InternalID,
		FeeID:           txCtx.FeeID,
		ContractID:      txCtx.ContractId,
	}

	recipientId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, "recipient")
	if err != nil {
		return nil, err
	}
	transfer.ToID = recipientId

	senderId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, "sender")
	if err != nil {
		return nil, err
	}
	transfer.FromID = senderId

	transfer.Amount, err = parseTransferDecimal(calldata, "amount")
	if err != nil {
		return nil, err
	}

	return []storage.Transfer{transfer}, nil
}

func (parser TransferParser) parseTransferFromERC721Calldata(ctx context.Context, txCtx data.TxContext, calldata map[string]any) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		InternalID:      txCtx.InternalID,
		FeeID:           txCtx.FeeID,
		Amount:          decimal.NewFromInt(1),
		ContractID:      txCtx.ContractId,
	}

	recipientId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, "to")
	if err != nil {
		return nil, err
	}
	transfer.ToID = recipientId

	senderId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, "from_")
	if err != nil {
		return nil, err
	}
	transfer.FromID = senderId

	transfer.TokenID, err = parseTransferDecimal(calldata, "tokenId")
	if err != nil {
		return nil, err
	}

	return []storage.Transfer{transfer}, nil
}

func (parser TransferParser) parseMintErc20(ctx context.Context, txCtx data.TxContext, calldata map[string]any) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		InternalID:      txCtx.InternalID,
		FeeID:           txCtx.FeeID,
		ContractID:      txCtx.ContractId,
	}

	amount, err := parseTransferDecimal(calldata, "amount")
	if err != nil {
		return nil, err
	}
	transfer.Amount = amount

	for _, key := range []string{
		"to", "recipient",
	} {
		if _, ok := calldata[key]; ok {
			recipientId, err := parseTransferAddress(ctx, parser.resolver, transfer.Height, calldata, key)
			if err != nil {
				return nil, err
			}
			transfer.ToID = recipientId
			return []storage.Transfer{transfer}, nil
		}

	}

	return nil, nil
}
