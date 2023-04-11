package grpc

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
)

type subscriptionFilters struct {
	declare       *storage.DeclareFilter
	deploy        *storage.DeployFilter
	deployAccount *storage.DeployAccountFilter
	event         *storage.EventFilter
	fee           *storage.FeeFilter
	internal      *storage.InternalFilter
	invoke        *storage.InvokeFilter
	l1Handler     *storage.L1HandlerFilter
	message       *storage.MessageFilter
	storageDiff   *storage.StorageDiffFilter
	tokenBalance  *storage.TokenBalanceFilter
	transfer      *storage.TransferFilter
}

func newSubscriptionFilters(req *pb.SubscribeRequest) subscriptionFilters {
	return subscriptionFilters{
		declare:       declareFilter(req.GetDeclares()),
		deploy:        deployFilter(req.GetDeploys()),
		deployAccount: deployAccountFilter(req.GetDeployAccounts()),
		event:         eventFilter(req.GetEvents()),
		fee:           feeFilter(req.GetFees()),
		internal:      internalFilter(req.GetInternals()),
		invoke:        invokeFilter(req.GetInvokes()),
		l1Handler:     l1HandlerFilter(req.GetL1Handlers()),
		message:       messageFilter(req.GetMsgs()),
		storageDiff:   storageDiffFilter(req.GetStorageDiffs()),
		tokenBalance:  tokenBalanceFilter(req.GetTokenBalances()),
		transfer:      transferFilter(req.GetTransfers()),
	}
}

func declareFilter(fltr *pb.DeclareFilters) *storage.DeclareFilter {
	if fltr == nil {
		return nil
	}
	return &storage.DeclareFilter{
		ID:      integerFilter(fltr.Id),
		Height:  integerFilter(fltr.Height),
		Time:    timeFilter(fltr.Time),
		Version: enumFilter(fltr.Version),
		Status:  enumFilter(fltr.Status),
	}
}

func deployFilter(fltr *pb.DeployFilters) *storage.DeployFilter {
	if fltr == nil {
		return nil
	}
	return &storage.DeployFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Class:          bytesFilter(fltr.Class),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func deployAccountFilter(fltr *pb.DeployAccountFilters) *storage.DeployAccountFilter {
	if fltr == nil {
		return nil
	}
	return &storage.DeployAccountFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Class:          bytesFilter(fltr.Class),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func eventFilter(fltr *pb.EventFilter) *storage.EventFilter {
	if fltr == nil {
		return nil
	}
	return &storage.EventFilter{
		ID:         integerFilter(fltr.Id),
		Height:     integerFilter(fltr.Height),
		Time:       timeFilter(fltr.Time),
		Contract:   bytesFilter(fltr.Contract),
		From:       bytesFilter(fltr.From),
		Name:       stringFilter(fltr.Name),
		ParsedData: fltr.GetParsedData(),
	}
}

func feeFilter(fltr *pb.FeeFilter) *storage.FeeFilter {
	if fltr == nil {
		return nil
	}
	return &storage.FeeFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Contract:       bytesFilter(fltr.Contract),
		Caller:         bytesFilter(fltr.Caller),
		Class:          bytesFilter(fltr.Class),
		Selector:       equalityFilter(fltr.Selector),
		Entrypoint:     stringFilter(fltr.Entrypoint),
		EntrypointType: enumFilter(fltr.EntrypointType),
		CallType:       enumFilter(fltr.CallType),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func internalFilter(fltr *pb.InternalFilter) *storage.InternalFilter {
	if fltr == nil {
		return nil
	}
	return &storage.InternalFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Contract:       bytesFilter(fltr.Contract),
		Caller:         bytesFilter(fltr.Caller),
		Class:          bytesFilter(fltr.Class),
		Selector:       equalityFilter(fltr.Selector),
		Entrypoint:     stringFilter(fltr.Entrypoint),
		EntrypointType: enumFilter(fltr.EntrypointType),
		CallType:       enumFilter(fltr.CallType),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func invokeFilter(fltr *pb.InvokeFilters) *storage.InvokeFilter {
	if fltr == nil {
		return nil
	}
	return &storage.InvokeFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Version:        enumFilter(fltr.Version),
		Contract:       bytesFilter(fltr.Contract),
		Selector:       equalityFilter(fltr.Selector),
		Entrypoint:     stringFilter(fltr.Entrypoint),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func l1HandlerFilter(fltr *pb.L1HandlerFilter) *storage.L1HandlerFilter {
	if fltr == nil {
		return nil
	}
	return &storage.L1HandlerFilter{
		ID:             integerFilter(fltr.Id),
		Height:         integerFilter(fltr.Height),
		Time:           timeFilter(fltr.Time),
		Status:         enumFilter(fltr.Status),
		Contract:       bytesFilter(fltr.Contract),
		Selector:       equalityFilter(fltr.Selector),
		Entrypoint:     stringFilter(fltr.Entrypoint),
		ParsedCalldata: fltr.GetParsedCalldata(),
	}
}

func messageFilter(fltr *pb.MessageFilter) *storage.MessageFilter {
	if fltr == nil {
		return nil
	}
	return &storage.MessageFilter{
		ID:       integerFilter(fltr.Id),
		Height:   integerFilter(fltr.Height),
		Time:     timeFilter(fltr.Time),
		Contract: bytesFilter(fltr.Contract),
		From:     bytesFilter(fltr.From),
		To:       bytesFilter(fltr.To),
		Selector: equalityFilter(fltr.Selector),
	}
}

func storageDiffFilter(fltr *pb.StorageDiffFilter) *storage.StorageDiffFilter {
	if fltr == nil {
		return nil
	}
	return &storage.StorageDiffFilter{
		ID:       integerFilter(fltr.Id),
		Height:   integerFilter(fltr.Height),
		Contract: bytesFilter(fltr.Contract),
		Key:      equalityFilter(fltr.Key),
	}
}

func tokenBalanceFilter(fltr *pb.TokenBalanceFilter) *storage.TokenBalanceFilter {
	if fltr == nil {
		return nil
	}
	return &storage.TokenBalanceFilter{
		Owner:    bytesFilter(fltr.Owner),
		Contract: bytesFilter(fltr.Contract),
		TokenId:  stringFilter(fltr.TokenId),
	}
}

func transferFilter(fltr *pb.TransferFilter) *storage.TransferFilter {
	if fltr == nil {
		return nil
	}
	return &storage.TransferFilter{
		ID:       integerFilter(fltr.Id),
		Height:   integerFilter(fltr.Height),
		Time:     timeFilter(fltr.Time),
		Contract: bytesFilter(fltr.Contract),
		From:     bytesFilter(fltr.From),
		To:       bytesFilter(fltr.To),
		TokenId:  stringFilter(fltr.TokenId),
	}
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
