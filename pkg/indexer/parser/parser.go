package parser

import (
	"context"
	"math/big"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Result -
type Result struct {
	Addresses map[string]*storage.Address
	Block     storage.Block
	Classes   map[string]*storage.Class
}

// Parser -
type Parser struct {
	receiver     *receiver.Receiver
	cache        *cache.Cache
	idGenerator  *IdGenerator
	storageDiffs storage.IStorageDiff

	addresses map[string]*storage.Address
	classes   map[string]*storage.Class
}

// New -
func New(
	receiver *receiver.Receiver,
	cache *cache.Cache,
	idGenerator *IdGenerator,
	storageDiffs storage.IStorageDiff,
) *Parser {
	return &Parser{
		receiver:     receiver,
		cache:        cache,
		idGenerator:  idGenerator,
		storageDiffs: storageDiffs,

		addresses: make(map[string]*storage.Address),
		classes:   make(map[string]*storage.Class),
	}
}

func (parser *Parser) addAddress(address *storage.Address) {
	if len(address.Hash) == 0 {
		return
	}
	key := encoding.EncodeHex(address.Hash)
	if _, ok := parser.addresses[key]; !ok {
		parser.addresses[key] = address
	}
}

func (parser *Parser) addClass(class *storage.Class) {
	if len(class.Hash) == 0 {
		return
	}
	key := encoding.EncodeHex(class.Hash)
	if _, ok := parser.classes[key]; !ok {
		parser.classes[key] = class
	}
}

// Parse -
func (parser *Parser) Parse(ctx context.Context, result receiver.Result) (Result, error) {
	block := storage.Block{
		ID:               result.Block.BlockNumber + 1,
		Height:           result.Block.BlockNumber,
		Time:             result.Block.Timestamp,
		Hash:             encoding.MustDecodeHex(result.Block.BlockHash),
		ParentHash:       encoding.MustDecodeHex(result.Block.ParentHash),
		NewRoot:          encoding.MustDecodeHex(result.Block.NewRoot),
		SequencerAddress: encoding.MustDecodeHex(result.Block.SequencerAddress),
		Status:           storage.NewStatus(result.Block.Status),
		TxCount:          len(result.Block.Transactions),

		InvokeV0:      make([]storage.InvokeV0, 0),
		InvokeV1:      make([]storage.InvokeV1, 0),
		Declare:       make([]storage.Declare, 0),
		Deploy:        make([]storage.Deploy, 0),
		DeployAccount: make([]storage.DeployAccount, 0),
		L1Handler:     make([]storage.L1Handler, 0),
	}

	if err := parser.entitiesFromStateUpdate(ctx, &block, result.StateUpdate); err != nil {
		return Result{}, errors.Wrap(err, "state update parsing")
	}

	for i := range result.Block.Transactions {
		switch typed := result.Block.Transactions[i].Body.(type) {
		case *data.InvokeV0:
			invoke, err := parser.getInvokeV0(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "invoke v0")
			}
			block.InvokeV0 = append(block.InvokeV0, invoke)
		case *data.InvokeV1:
			invoke, err := parser.getInvokeV1(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "invoke v1")
			}
			block.InvokeV1 = append(block.InvokeV1, invoke)
		case *data.Declare:
			tx, err := parser.getDeclare(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "declare")
			}
			block.Declare = append(block.Declare, tx)
		case *data.Deploy:
			tx, err := parser.getDeploy(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "deploy")
			}
			block.Deploy = append(block.Deploy, tx)
		case *data.DeployAccount:
			tx, err := parser.getDeployAccount(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "deploy account")
			}
			block.DeployAccount = append(block.DeployAccount, tx)
		case *data.L1Handler:
			tx, err := parser.getL1Handler(ctx, typed, block, result.Trace.Traces[i])
			if err != nil {
				return Result{}, errors.Wrap(err, "l1 handler")
			}
			block.L1Handler = append(block.L1Handler, tx)
		default:
			return Result{}, errors.Errorf("unknown transaction type: %s", result.Block.Transactions[i].Type)
		}
	}

	block.InvokeV0Count = len(block.InvokeV0)
	block.InvokeV1Count = len(block.InvokeV1)
	block.DeclareCount = len(block.Declare)
	block.DeployCount = len(block.Deploy)
	block.DeployAccountCount = len(block.DeployAccount)
	block.L1HandlerCount = len(block.L1Handler)

	return Result{
		Block:     block,
		Addresses: parser.addresses,
		Classes:   parser.classes,
	}, nil
}

func decimalFromHex(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	i, _ := new(big.Int).SetString(s, 0)
	return decimal.NewFromBigInt(i, 0)
}

func (parser *Parser) receiveClass(ctx context.Context, class *storage.Class) error {
	rawClass, err := parser.receiver.GetClass(ctx, encoding.EncodeHex(class.Hash))
	if err != nil {
		return err
	}

	class.Abi = storage.Bytes(rawClass.RawAbi)

	a, err := rawClass.GetAbi()
	if err != nil {
		return err
	}
	interfaces, err := starknet.FindInterfaces(a)
	if err != nil {
		return err
	}
	class.Type = storage.NewClassType(interfaces...)

	parser.cache.SetAbiByClassHash(*class, a)
	parser.addClass(class)

	return nil
}

func (parser *Parser) findAddress(ctx context.Context, address *storage.Address) error {
	if value, ok := parser.addresses[encoding.EncodeHex(address.Hash)]; ok {
		address.ID = value.ID
		address.ClassID = value.ClassID
		return nil
	}
	generated, err := parser.idGenerator.SetAddressId(ctx, address)
	if err != nil {
		return err
	}
	if generated {
		parser.addAddress(address)
	}
	return nil
}

func (parser *Parser) findClass(ctx context.Context, class *storage.Class) error {
	if value, ok := parser.classes[encoding.EncodeHex(class.Hash)]; ok {
		class.ID = value.ID
		class.Abi = value.Abi
		class.Type = value.Type
		return nil
	}
	generated, err := parser.idGenerator.SetClassId(ctx, class)
	if err != nil {
		return err
	}
	if generated {
		parser.addClass(class)
	}
	return nil
}
