package v0

import (
	"bytes"
	"context"
	"errors"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	parserData "github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
)

// ParseDeploy -
func (parser Parser) ParseDeploy(ctx context.Context, raw *data.Deploy, block storage.Block, trace sequencer.Trace, receipts sequencer.Receipt) (storage.Deploy, *storage.Fee, error) {
	tx := storage.Deploy{
		ID:                  parser.Resolver.NextTxId(),
		Height:              block.Height,
		Time:                block.Time,
		Status:              block.Status,
		Hash:                trace.TransactionHash.Bytes(),
		ContractAddressSalt: encoding.MustDecodeHex(raw.ContractAddressSalt),
		ConstructorCalldata: raw.ConstructorCalldata,
	}

	if class, err := parser.Resolver.FindClassByHash(ctx, raw.ClassHash); err != nil {
		return tx, nil, err
	} else if class != nil {
		tx.Class = *class
		tx.ClassID = class.ID
	}

	if address, err := parser.Resolver.FindAddressByHash(ctx, raw.ContractAddress); err != nil {
		return tx, nil, err
	} else if address != nil {
		tx.ContractID = address.ID
		tx.Contract = *address
		tx.Contract.Height = tx.Height
		if tx.ClassID != 0 {
			tx.Contract.ClassID = &tx.Class.ID
			tx.Contract.Class = tx.Class
		}
	}

	classAbi, err := parser.Cache.GetAbiByClassHash(ctx, tx.Class.Hash)
	if err != nil {
		return tx, nil, err
	}

	if len(tx.ConstructorCalldata) > 0 {
		tx.ParsedCalldata, err = decode.CalldataForConstructor(classAbi, tx.ConstructorCalldata)
		if err != nil {
			if !errors.Is(err, abi.ErrNoLenField) {
				return tx, nil, err
			}
		}
	}

	parser.Cache.SetAbiByAddress(tx.Class, tx.Contract.Hash)

	var proxyId uint64
	if tx.Class.Type.Is(storage.ClassTypeProxy) {
		proxyId = tx.ContractID
	}

	txCtx := parserData.NewTxContextFromDeploy(tx, proxyId)

	if trace.FunctionInvocation != nil {
		tx.Events, err = parseEvents(ctx, parser.EventParser, txCtx, classAbi, trace.FunctionInvocation.Events)
		if err != nil {
			return tx, nil, err
		}
		tx.Transfers, err = parser.TransferParser.ParseEvents(ctx, txCtx, tx.Contract, tx.Events)
		if err != nil {
			return tx, nil, err
		}

		tx.Messages, err = parseMessages(ctx, parser.MessageParser, txCtx, trace.FunctionInvocation.Messages)
		if err != nil {
			return tx, nil, err
		}

		tx.Internals, err = parseInternals(ctx, parser.InternalTxParser, txCtx, trace.FunctionInvocation.InternalCalls)
		if err != nil {
			return tx, nil, err
		}
	}

	if trace.FeeTransferInvocation != nil {
		fee, err := parser.FeeParser.ParseInvocation(ctx, txCtx, *trace.FeeTransferInvocation)
		if err != nil {
			return tx, nil, nil
		}
		if fee != nil {
			return tx, fee, nil
		}
	} else {
		transfer, err := parser.FeeParser.ParseActualFee(ctx, txCtx, receipts.ActualFee)
		if err != nil {
			return tx, nil, nil
		}
		if transfer != nil {
			tx.Transfers = append(tx.Transfers, *transfer)
		}
	}

	if token := createBridgedToken(ctx, block, tx.Contract); token != nil {
		tx.ERC20 = token
	} else if tx.Class.Type.OneOf(storage.ClassTypeERC20, storage.ClassTypeERC721, storage.ClassTypeERC1155) {
		token, err := parser.TokenParser.Parse(ctx, txCtx, tx.Contract, tx.Class.Type, tx.ParsedCalldata)
		if err != nil {
			return tx, nil, err
		}
		tx.ERC20 = token.ERC20
		tx.ERC721 = token.ERC721
		tx.ERC1155 = token.ERC1155
	}

	return tx, nil, nil
}

func createBridgedToken(ctx context.Context, block storage.Block, contract storage.Address) *storage.ERC20 {
	tokens := starknet.BridgedTokens()

	for i := range tokens {
		if !bytes.Equal(tokens[i].L2TokenAddress.Bytes(), contract.Hash) {
			continue
		}

		return &storage.ERC20{
			DeployHeight: block.Height,
			DeployTime:   block.Time,
			ContractID:   contract.ID,
			Name:         tokens[i].Name,
			Symbol:       tokens[i].Symbol,
			Decimals:     tokens[i].Decimals,
		}
	}

	return nil
}
