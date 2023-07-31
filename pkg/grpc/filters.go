package grpc

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

type subscriptionFilters struct {
	declare       []storage.DeclareFilter
	deploy        []storage.DeployFilter
	deployAccount []storage.DeployAccountFilter
	event         []storage.EventFilter
	fee           []storage.FeeFilter
	internal      []storage.InternalFilter
	invoke        []storage.InvokeFilter
	l1Handler     []storage.L1HandlerFilter
	message       []storage.MessageFilter
	storageDiff   []storage.StorageDiffFilter
	tokenBalance  []storage.TokenBalanceFilter
	transfer      []storage.TransferFilter
	address       []storage.AddressFilter
	tokens        []storage.TokenFilter
}

func newSubscriptionFilters(ctx context.Context, req *pb.SubscribeRequest, db postgres.Storage) (subscriptionFilters, error) {
	event, err := eventFilter(ctx, db.Address, req.GetEvents())
	if err != nil {
		return subscriptionFilters{}, err
	}
	return subscriptionFilters{
		declare:       declareFilter(req.GetDeclares()),
		deploy:        deployFilter(req.GetDeploys()),
		deployAccount: deployAccountFilter(req.GetDeployAccounts()),
		event:         event,
		fee:           feeFilter(req.GetFees()),
		internal:      internalFilter(req.GetInternals()),
		invoke:        invokeFilter(req.GetInvokes()),
		l1Handler:     l1HandlerFilter(req.GetL1Handlers()),
		message:       messageFilter(req.GetMsgs()),
		storageDiff:   storageDiffFilter(req.GetStorageDiffs()),
		tokenBalance:  tokenBalanceFilter(req.GetTokenBalances()),
		transfer:      transferFilter(req.GetTransfers()),
		address:       addressFilter(req.GetAddresses()),
		tokens:        tokenFilter(req.GetTokens()),
	}, nil
}

func addressFilter(fltr []*pb.AddressFilter) []storage.AddressFilter {
	if len(fltr) == 0 {
		return nil
	}
	result := make([]storage.AddressFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.AddressFilter{
			ID:           integerFilter(fltr[i].Id),
			OnlyStarknet: fltr[i].OnlyStarknet,
			Height:       integerFilter(fltr[i].Height),
		}
	}
	return result
}

func declareFilter(fltr []*pb.DeclareFilters) []storage.DeclareFilter {
	if len(fltr) == 0 {
		return nil
	}
	result := make([]storage.DeclareFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.DeclareFilter{
			ID:      integerFilter(fltr[i].Id),
			Height:  integerFilter(fltr[i].Height),
			Time:    timeFilter(fltr[i].Time),
			Version: enumFilter(fltr[i].Version),
			Status:  enumFilter(fltr[i].Status),
		}
	}
	return result
}

func deployFilter(fltr []*pb.DeployFilters) []storage.DeployFilter {
	if len(fltr) == 0 {
		return nil
	}
	result := make([]storage.DeployFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.DeployFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Class:          bytesFilter(fltr[i].Class),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func deployAccountFilter(fltr []*pb.DeployAccountFilters) []storage.DeployAccountFilter {
	if len(fltr) == 0 {
		return nil
	}
	result := make([]storage.DeployAccountFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.DeployAccountFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Class:          bytesFilter(fltr[i].Class),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func eventFilter(ctx context.Context, address storage.IAddress, fltr []*pb.EventFilter) ([]storage.EventFilter, error) {
	if len(fltr) == 0 {
		return nil, nil
	}

	result := make([]storage.EventFilter, len(fltr))
	for i := range fltr {
		contractFilter, err := idFilter(ctx, address, fltr[i].Contract)
		if err != nil {
			return nil, err
		}
		fromFilter, err := idFilter(ctx, address, fltr[i].From)
		if err != nil {
			return nil, err
		}
		result[i] = storage.EventFilter{
			ID:         integerFilter(fltr[i].Id),
			Height:     integerFilter(fltr[i].Height),
			Time:       timeFilter(fltr[i].Time),
			Contract:   contractFilter,
			From:       fromFilter,
			Name:       stringFilter(fltr[i].Name),
			ParsedData: fltr[i].GetParsedData(),
		}
	}

	return result, nil
}

func feeFilter(fltr []*pb.FeeFilter) []storage.FeeFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.FeeFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.FeeFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Contract:       bytesFilter(fltr[i].Contract),
			Caller:         bytesFilter(fltr[i].Caller),
			Class:          bytesFilter(fltr[i].Class),
			Selector:       equalityFilter(fltr[i].Selector),
			Entrypoint:     stringFilter(fltr[i].Entrypoint),
			EntrypointType: enumFilter(fltr[i].EntrypointType),
			CallType:       enumFilter(fltr[i].CallType),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func internalFilter(fltr []*pb.InternalFilter) []storage.InternalFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.InternalFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.InternalFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Contract:       bytesFilter(fltr[i].Contract),
			Caller:         bytesFilter(fltr[i].Caller),
			Class:          bytesFilter(fltr[i].Class),
			Selector:       equalityFilter(fltr[i].Selector),
			Entrypoint:     stringFilter(fltr[i].Entrypoint),
			EntrypointType: enumFilter(fltr[i].EntrypointType),
			CallType:       enumFilter(fltr[i].CallType),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func invokeFilter(fltr []*pb.InvokeFilters) []storage.InvokeFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.InvokeFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.InvokeFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Version:        enumFilter(fltr[i].Version),
			Contract:       bytesFilter(fltr[i].Contract),
			Selector:       equalityFilter(fltr[i].Selector),
			Entrypoint:     stringFilter(fltr[i].Entrypoint),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func l1HandlerFilter(fltr []*pb.L1HandlerFilter) []storage.L1HandlerFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.L1HandlerFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.L1HandlerFilter{
			ID:             integerFilter(fltr[i].Id),
			Height:         integerFilter(fltr[i].Height),
			Time:           timeFilter(fltr[i].Time),
			Status:         enumFilter(fltr[i].Status),
			Contract:       bytesFilter(fltr[i].Contract),
			Selector:       equalityFilter(fltr[i].Selector),
			Entrypoint:     stringFilter(fltr[i].Entrypoint),
			ParsedCalldata: fltr[i].GetParsedCalldata(),
		}
	}
	return result
}

func messageFilter(fltr []*pb.MessageFilter) []storage.MessageFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.MessageFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.MessageFilter{
			ID:       integerFilter(fltr[i].Id),
			Height:   integerFilter(fltr[i].Height),
			Time:     timeFilter(fltr[i].Time),
			Contract: bytesFilter(fltr[i].Contract),
			From:     bytesFilter(fltr[i].From),
			To:       bytesFilter(fltr[i].To),
			Selector: equalityFilter(fltr[i].Selector),
		}
	}
	return result
}

func storageDiffFilter(fltr []*pb.StorageDiffFilter) []storage.StorageDiffFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.StorageDiffFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.StorageDiffFilter{
			ID:       integerFilter(fltr[i].Id),
			Height:   integerFilter(fltr[i].Height),
			Contract: bytesFilter(fltr[i].Contract),
			Key:      equalityFilter(fltr[i].Key),
		}
	}
	return result
}

func tokenFilter(fltr []*pb.TokenFilter) []storage.TokenFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.TokenFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.TokenFilter{
			ID:       integerFilter(fltr[i].Id),
			Type:     enumStringFilter(fltr[i].Type),
			TokenId:  stringFilter(fltr[i].TokenId),
			Contract: bytesFilter(fltr[i].Contract),
		}
	}
	return result
}

func tokenBalanceFilter(fltr []*pb.TokenBalanceFilter) []storage.TokenBalanceFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.TokenBalanceFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.TokenBalanceFilter{
			Owner:    bytesFilter(fltr[i].Owner),
			Contract: bytesFilter(fltr[i].Contract),
			TokenId:  stringFilter(fltr[i].TokenId),
		}
	}
	return result
}

func transferFilter(fltr []*pb.TransferFilter) []storage.TransferFilter {
	if len(fltr) == 0 {
		return nil
	}

	result := make([]storage.TransferFilter, len(fltr))
	for i := range fltr {
		result[i] = storage.TransferFilter{
			ID:       integerFilter(fltr[i].Id),
			Height:   integerFilter(fltr[i].Height),
			Time:     timeFilter(fltr[i].Time),
			Contract: bytesFilter(fltr[i].Contract),
			From:     bytesFilter(fltr[i].From),
			To:       bytesFilter(fltr[i].To),
			TokenId:  stringFilter(fltr[i].TokenId),
		}
	}
	return result
}

func integerFilter(fltr *pb.IntegerFilter) (result storage.IntegerFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	result.Gt = fltr.GetGt()
	result.Gte = fltr.GetGte()
	result.Lt = fltr.GetLt()
	result.Lte = fltr.GetLte()
	result.Neq = fltr.GetNeq()
	result.Between = betweenFilter(fltr.GetBetween())
	return
}

func betweenFilter(fltr *pb.BetweenInteger) *storage.BetweenFilter {
	if fltr == nil {
		return nil
	}

	result := new(storage.BetweenFilter)
	result.From = fltr.GetFrom()
	result.To = fltr.GetTo()
	return result
}

func timeFilter(fltr *pb.TimeFilter) (result storage.TimeFilter) {
	if fltr == nil {
		return
	}

	result.Gt = fltr.GetGt()
	result.Gte = fltr.GetGte()
	result.Lt = fltr.GetLt()
	result.Lte = fltr.GetLte()
	result.Between = betweenFilter(fltr.GetBetween())

	return
}

func enumFilter(fltr *pb.EnumFilter) (result storage.EnumFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	result.Neq = fltr.GetNeq()
	if arr := fltr.GetIn(); arr != nil {
		result.In = arr.GetArr()
	}
	if arr := fltr.GetNotin(); arr != nil {
		result.Notin = arr.GetArr()
	}

	return
}

func enumStringFilter(fltr *pb.EnumStringFilter) (result storage.EnumStringFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	result.Neq = fltr.GetNeq()
	if arr := fltr.GetIn(); arr != nil {
		result.In = arr.GetArr()
	}
	if arr := fltr.GetNotin(); arr != nil {
		result.Notin = arr.GetArr()
	}

	return
}

func bytesFilter(fltr *pb.BytesFilter) (result storage.BytesFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	if arr := fltr.GetIn(); arr != nil {
		result.In = arr.GetArr()
	}

	return
}

func equalityFilter(fltr *pb.EqualityFilter) (result storage.EqualityFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	result.Neq = fltr.GetNeq()

	return
}

func stringFilter(fltr *pb.StringFilter) (result storage.StringFilter) {
	if fltr == nil {
		return
	}

	result.Eq = fltr.GetEq()
	if arr := fltr.GetIn(); arr != nil {
		result.In = arr.GetArr()
	}

	return
}

func idFilter(ctx context.Context, address storage.IAddress, fltr *pb.BytesFilter) (result storage.IdFilter, err error) {
	if fltr == nil {
		return
	}

	switch {
	case len(fltr.GetEq()) > 0:
		a, err := address.GetByHash(ctx, fltr.GetEq())
		if err != nil {
			return result, err
		}
		result.Eq = a.ID
	case fltr.GetIn() != nil && len(fltr.GetIn().Arr) > 0:
		ids, err := address.GetIdsByHash(ctx, fltr.GetIn().Arr)
		if err != nil {
			return result, err
		}
		result.In = ids
	}

	return
}
