package v0

import (
	"bytes"
	"context"
	"strconv"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/decode"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/data"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/resolver"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	errInvalidTransfer = errors.New("invalid transfer")
)

// EventParser -
type EventParser struct {
	cache    *cache.Cache
	resolver resolver.Resolver
}

// NewEventParser -
func NewEventParser(
	cache *cache.Cache,
	resolver resolver.Resolver,
) EventParser {
	return EventParser{
		cache:    cache,
		resolver: resolver,
	}
}

// Parse -
func (parser EventParser) Parse(ctx context.Context, txCtx data.TxContext, contractAbi abi.Abi, event starknetData.Event) (storage.Event, error) {
	model := storage.Event{
		ID:              parser.resolver.NextEventId(),
		Height:          txCtx.Height,
		Time:            txCtx.Time,
		Order:           event.Order,
		Data:            make([]string, len(event.Data)),
		Keys:            make([]string, len(event.Keys)),
		ContractID:      txCtx.ContractId,
		Contract:        txCtx.Contract,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
	}
	for i := range event.Data {
		model.Data[i] = event.Data[i].String()
	}
	for i := range event.Keys {
		model.Keys[i] = event.Keys[i].String()
	}

	if address, err := parser.resolver.FindAddressByHash(ctx, starknetData.Felt(event.FromAddress)); err != nil {
		return model, err
	} else if address != nil {
		model.FromID = address.ID
		model.From = *address
	}

	if len(contractAbi.Events) > 0 {
		parsed, name, err := decode.Event(contractAbi, model.Keys, model.Data)
		if err != nil {
			return model, err
		}
		model.ParsedData = parsed
		model.Name = name
	}

	return model, nil
}

func upgraded(params map[string]any, height uint64) ([]data.ProxyUpgrade, error) {
	value, ok := params["implementation"]
	if !ok {
		return nil, nil
	}
	implementation, ok := value.(string)
	if !ok {
		return nil, nil
	}
	return []data.ProxyUpgrade{
		{
			Address: starknetData.Felt(implementation).Bytes(),
			Action:  storage.ProxyActionUpdate,
			Height:  height,
		},
	}, nil
}

func accountUpgraded(params map[string]any, height uint64) ([]data.ProxyUpgrade, error) {
	value, ok := params["new_implementation"]
	if !ok {
		return nil, nil
	}
	implementation, ok := value.(string)
	if !ok {
		return nil, nil
	}
	return []data.ProxyUpgrade{
		{
			Address: starknetData.Felt(implementation).Bytes(),
			Action:  storage.ProxyActionUpdate,
			Height:  height,
		},
	}, nil
}

func implementationUpgraded(params map[string]any, height uint64) ([]data.ProxyUpgrade, error) {
	value, ok := params["implementation_hash"]
	if !ok {
		return nil, nil
	}
	implementation, ok := value.(string)
	if !ok {
		return nil, nil
	}
	return []data.ProxyUpgrade{
		{
			Address: starknetData.Felt(implementation).Bytes(),
			Action:  storage.ProxyActionUpdate,
			Height:  height,
		},
	}, nil
}

func moduleFunctionChange(params map[string]any, height uint64) ([]data.ProxyUpgrade, error) {
	upgrade := new(data.ProxyUpgrade)
	actionsValue, ok := params["actions"]
	if !ok {
		return nil, nil
	}
	actions, ok := actionsValue.([]any)
	if !ok {
		return nil, nil
	}

	upgrades := make([]data.ProxyUpgrade, 0)
	for i := range actions {
		action, ok := actions[i].(map[string]any)
		if !ok {
			continue
		}

		value, ok := action["module_address"]
		if !ok {
			continue
		}
		address, ok := value.(string)
		if !ok {
			continue
		}
		upgrade.Address = starknetData.Felt(address).Bytes()

		value, ok = action["selector"]
		if !ok {
			continue
		}
		selector, ok := value.(string)
		if !ok {
			continue
		}
		upgrade.Selector = starknetData.Felt(selector).Bytes()

		value, ok = action["action"]
		if !ok {
			return nil, nil
		}
		actionType, ok := value.(string)
		if !ok {
			continue
		}
		iAction, err := strconv.ParseInt(actionType, 0, 32)
		if err != nil {
			return nil, err
		}
		upgrade.Action = storage.ProxyAction(iAction)
		upgrade.IsModule = true
		upgrade.Height = height
		upgrades = append(upgrades, *upgrade)
	}

	return upgrades, nil
}

func depositHandled(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contract storage.Address, event storage.Event) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          event.Height,
		Time:            event.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
	}

	bridged := starknet.BridgedTokens()
	for i := range bridged {
		if bytes.Equal(contract.Hash, bridged[i].L2BridgeAddress.Bytes()) {
			address, err := resolver.FindAddressByHash(ctx, bridged[i].L2TokenAddress)
			if err != nil {
				return nil, err
			}
			transfer.ContractID = address.ID
			break
		}
	}

	transfer.Token = storage.Token{
		FirstHeight: event.Height,
		ContractId:  transfer.ContractID,
		Type:        storage.TokenTypeERC20,
		TokenId:     decimal.Zero,
	}

	var err error

	transfer.ToID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, "account")
	if err != nil {
		return nil, err
	}

	if transfer.ToID == 0 {
		return nil, errInvalidTransfer
	}

	transfer.Amount, err = parseTransferDecimal(event.ParsedData, "amount")
	if err != nil {
		return nil, err
	}

	return []storage.Transfer{transfer}, nil
}

func withdrawInitiated(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contract storage.Address, event storage.Event) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          event.Height,
		Time:            event.Time,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
	}

	bridged := starknet.BridgedTokens()
	for i := range bridged {
		if bytes.Equal(contract.Hash, bridged[i].L2BridgeAddress.Bytes()) {
			address, err := resolver.FindAddressByHash(ctx, bridged[i].L2TokenAddress)
			if err != nil {
				return nil, err
			}
			transfer.ContractID = address.ID
			break
		}
	}

	transfer.Token = storage.Token{
		FirstHeight: event.Height,
		ContractId:  transfer.ContractID,
		Type:        storage.TokenTypeERC20,
		TokenId:     decimal.Zero,
	}

	var err error

	transfer.FromID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, "caller_address")
	if err != nil {
		return nil, err
	}

	if transfer.FromID == 0 {
		return nil, errInvalidTransfer
	}

	transfer.Amount, err = parseTransferDecimal(event.ParsedData, "amount")
	if err != nil {
		return nil, err
	}

	return []storage.Transfer{transfer}, nil
}

func transfer(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contractId uint64, event storage.Event) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          event.Height,
		Time:            event.Time,
		ContractID:      contractId,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
		Token: storage.Token{
			FirstHeight: event.Height,
			ContractId:  contractId,
			Type:        storage.TokenTypeERC20,
			TokenId:     decimal.Zero,
		},
	}

	var err error

	for _, key := range []string{
		"from_", "sender", "from_address",
	} {
		transfer.FromID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.FromID > 0 {
			break
		}
	}
	for _, key := range []string{
		"to", "recipient", "to_address",
	} {
		transfer.ToID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.ToID > 0 {
			break
		}
	}

	if transfer.FromID == 0 && transfer.ToID == 0 {
		return nil, errInvalidTransfer
	}

	transfer.Amount, err = parseTransferDecimal(event.ParsedData, "value")
	if err != nil {
		return nil, err
	}

	return []storage.Transfer{transfer}, nil
}

func transferERC721(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contractId uint64, event storage.Event) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          event.Height,
		Time:            event.Time,
		ContractID:      contractId,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
		Amount:          decimal.NewFromInt(1),
	}

	var err error

	for _, key := range []string{
		"from_", "sender",
	} {
		transfer.FromID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.FromID > 0 {
			break
		}
	}
	for _, key := range []string{
		"to", "recipient",
	} {
		transfer.ToID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.ToID > 0 {
			break
		}
	}

	if transfer.FromID == 0 && transfer.ToID == 0 {
		return nil, errInvalidTransfer
	}

	transfer.TokenID, err = parseTransferDecimal(event.ParsedData, "tokenId")
	if err != nil {
		return nil, err
	}

	transfer.Token = storage.Token{
		FirstHeight: txCtx.Height,
		ContractId:  contractId,
		TokenId:     transfer.TokenID,
		Type:        storage.TokenTypeERC721,
	}

	return []storage.Transfer{transfer}, nil
}

func transferSingle(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contractId uint64, event storage.Event) ([]storage.Transfer, error) {
	transfer := storage.Transfer{
		Height:          event.Height,
		Time:            event.Time,
		ContractID:      contractId,
		DeclareID:       txCtx.DeclareID,
		DeployID:        txCtx.DeployID,
		DeployAccountID: txCtx.DeployAccountID,
		InvokeID:        txCtx.InvokeID,
		L1HandlerID:     txCtx.L1HandlerID,
		FeeID:           txCtx.FeeID,
		InternalID:      txCtx.InternalID,
	}

	var err error

	for _, key := range []string{"from_", "_from"} {
		transfer.FromID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.FromID > 0 {
			break
		}
	}
	for _, key := range []string{"to", "_to"} {
		transfer.ToID, err = parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.ToID > 0 {
			break
		}
	}

	if transfer.FromID == 0 && transfer.ToID == 0 {
		return nil, errInvalidTransfer
	}

	for _, key := range []string{"value", "_value"} {
		transfer.Amount, err = parseTransferDecimal(event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.Amount.IsPositive() {
			break
		}
	}

	for _, key := range []string{"id", "_id"} {
		transfer.TokenID, err = parseTransferDecimal(event.ParsedData, key)
		if err != nil {
			return nil, err
		}
		if transfer.TokenID.IsPositive() {
			break
		}
	}

	transfer.Token = storage.Token{
		FirstHeight: txCtx.Height,
		ContractId:  contractId,
		TokenId:     transfer.TokenID,
		Type:        storage.TokenTypeERC1155,
	}

	return []storage.Transfer{transfer}, nil
}

func transferBatch(ctx context.Context, resolver resolver.Resolver, txCtx data.TxContext, contractId uint64, event storage.Event) ([]storage.Transfer, error) {
	fromId, err := parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, "from_")
	if err != nil {
		return nil, err
	}

	toId, err := parseTransferAddress(ctx, resolver, event.Height, event.ParsedData, "to")
	if err != nil {
		return nil, err
	}

	if fromId == 0 && toId == 0 {
		return nil, errInvalidTransfer
	}

	ids, err := parseTransferArrayDecimals(event, "ids")
	if err != nil {
		return nil, err
	}
	values, err := parseTransferArrayDecimals(event, "values")
	if err != nil {
		return nil, err
	}

	if len(ids) != len(values) {
		return nil, errInvalidTransfer
	}

	transfers := make([]storage.Transfer, len(ids))
	for i := range ids {
		transfers[i] = storage.Transfer{
			Height:          event.Height,
			Time:            event.Time,
			ContractID:      contractId,
			DeclareID:       txCtx.DeclareID,
			DeployID:        txCtx.DeployID,
			DeployAccountID: txCtx.DeployAccountID,
			InvokeID:        txCtx.InvokeID,
			L1HandlerID:     txCtx.L1HandlerID,
			FeeID:           txCtx.FeeID,
			InternalID:      txCtx.InternalID,
			FromID:          fromId,
			ToID:            toId,
			Amount:          values[i],
			TokenID:         ids[i],
			Token: storage.Token{
				FirstHeight: txCtx.Height,
				ContractId:  contractId,
				TokenId:     ids[i],
				Type:        storage.TokenTypeERC1155,
			},
		}
	}
	return transfers, nil
}

func parseTransferAddress(ctx context.Context, resolver resolver.Resolver, height uint64, data map[string]any, key string) (uint64, error) {
	if value, ok := data[key]; ok {
		if sValue, ok := value.(string); ok {
			address := storage.Address{
				Hash:   starknetData.Felt(sValue).Bytes(),
				Height: height,
			}
			if err := resolver.FindAddress(ctx, &address); err != nil {
				return 0, err
			}
			return address.ID, nil
		}
	}
	return 0, nil
}

func parseDecimalValue(value any) (decimal.Decimal, error) {
	switch typ := value.(type) {
	case string:
		return encoding.DecimalFromHex(typ), nil
	case map[string]any:
		lowValue, ok := typ["low"]
		if !ok {
			return decimal.Zero, nil
		}
		low, ok := lowValue.(string)
		if !ok {
			return decimal.Zero, nil
		}
		highValue, ok := typ["high"]
		if !ok {
			return decimal.Zero, nil
		}
		high, ok := highValue.(string)
		if !ok {
			return decimal.Zero, nil
		}
		uint256 := starknetData.NewUint256FromStrings(low, high)
		return uint256.Decimal()
	default:
		return decimal.Zero, nil
	}
}

func parseTransferDecimal(data map[string]any, key string) (decimal.Decimal, error) {
	if value, ok := data[key]; ok {
		return parseDecimalValue(value)
	}
	return decimal.Zero, nil
}

func parseTransferArrayDecimals(event storage.Event, key string) ([]decimal.Decimal, error) {
	var err error

	if value, ok := event.ParsedData[key]; ok {
		if arr, ok := value.([]any); ok {
			result := make([]decimal.Decimal, len(arr))
			for i := range arr {
				result[i], err = parseDecimalValue(arr[i])
				if err != nil {
					return nil, err
				}
			}
			return result, nil
		}
	}
	return nil, nil
}
