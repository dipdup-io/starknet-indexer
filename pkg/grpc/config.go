package grpc

import (
	"encoding/hex"
	"strings"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"gopkg.in/yaml.v3"
)

// ClientConfig -
type ClientConfig struct {
	ServerAddress string                  `yaml:"server_address" validate:"required"`
	Subscriptions map[string]Subscription `yaml:"subscriptions" validate:"omitempty"`
}

// Subscription -
type Subscription struct {
	Head                 bool                    `yaml:"head" validate:"omitempty"`
	InvokeFilters        []*InvokeFilters        `yaml:"invokes" validate:"omitempty"`
	DeclareFilters       []*DeclareFilters       `yaml:"declares" validate:"omitempty"`
	DeployFilters        []*DeployFilters        `yaml:"deploys" validate:"omitempty"`
	DeployAccountFilters []*DeployAccountFilters `yaml:"deploy_accounts" validate:"omitempty"`
	L1HandlerFilter      []*L1HandlerFilters     `yaml:"l1_handlers" validate:"omitempty"`
	InternalFilter       []*InternalFilters      `yaml:"internals" validate:"omitempty"`
	FeeFilter            []*FeeFilters           `yaml:"fees" validate:"omitempty"`
	EventFilter          []*EventFilter          `yaml:"events" validate:"omitempty"`
	MessageFilter        []*MessageFilter        `yaml:"messages" validate:"omitempty"`
	TransferFilter       []*TransferFilter       `yaml:"transfers" validate:"omitempty"`
	StorageDiffFilter    []*StorageDiffFilter    `yaml:"storage_diffs" validate:"omitempty"`
	TokenBalanceFilter   []*TokenBalanceFilter   `yaml:"token_balances" validate:"omitempty"`
	AddressFilter        []*AddressFilter        `yaml:"addresses" validate:"omitempty"`
	TokenFilter          []*TokenFilter          `yaml:"tokens" validate:"omitempty"`
}

// ToGrpcFilter -
func (f Subscription) ToGrpcFilter() *pb.SubscribeRequest {
	req := new(pb.SubscribeRequest)
	req.Head = f.Head

	if len(f.AddressFilter) > 0 {
		req.Addresses = make([]*pb.AddressFilter, len(f.AddressFilter))
		for i := range f.AddressFilter {
			req.Addresses[i] = f.AddressFilter[i].ToGrpcFilter()
		}
	}
	if len(f.EventFilter) > 0 {
		req.Events = make([]*pb.EventFilter, len(f.EventFilter))
		for i := range f.EventFilter {
			req.Events[i] = f.EventFilter[i].ToGrpcFilter()
		}
	}
	if len(f.InvokeFilters) > 0 {
		req.Invokes = make([]*pb.InvokeFilters, len(f.InvokeFilters))
		for i := range f.InvokeFilters {
			req.Invokes[i] = f.InvokeFilters[i].ToGrpcFilter()
		}
	}
	if len(f.DeclareFilters) > 0 {
		req.Declares = make([]*pb.DeclareFilters, len(f.DeclareFilters))
		for i := range f.DeclareFilters {
			req.Declares[i] = f.DeclareFilters[i].ToGrpcFilter()
		}
	}
	if len(f.DeployFilters) > 0 {
		req.Deploys = make([]*pb.DeployFilters, len(f.DeployFilters))
		for i := range f.DeployFilters {
			req.Deploys[i] = f.DeployFilters[i].ToGrpcFilter()
		}
	}
	if len(f.DeployAccountFilters) > 0 {
		req.DeployAccounts = make([]*pb.DeployAccountFilters, len(f.DeployAccountFilters))
		for i := range f.DeployAccountFilters {
			req.DeployAccounts[i] = f.DeployAccountFilters[i].ToGrpcFilter()
		}
	}
	if len(f.L1HandlerFilter) > 0 {
		req.L1Handlers = make([]*pb.L1HandlerFilter, len(f.L1HandlerFilter))
		for i := range f.L1HandlerFilter {
			req.L1Handlers[i] = f.L1HandlerFilter[i].ToGrpcFilter()
		}
	}
	if len(f.InternalFilter) > 0 {
		req.Internals = make([]*pb.InternalFilter, len(f.InternalFilter))
		for i := range f.InternalFilter {
			req.Internals[i] = f.InternalFilter[i].ToGrpcFilter()
		}
	}
	if len(f.FeeFilter) > 0 {
		req.Fees = make([]*pb.FeeFilter, len(f.FeeFilter))
		for i := range f.FeeFilter {
			req.Fees[i] = f.FeeFilter[i].ToGrpcFilter()
		}
	}
	if len(f.MessageFilter) > 0 {
		req.Msgs = make([]*pb.MessageFilter, len(f.MessageFilter))
		for i := range f.MessageFilter {
			req.Msgs[i] = f.MessageFilter[i].ToGrpcFilter()
		}
	}
	if len(f.TransferFilter) > 0 {
		req.Transfers = make([]*pb.TransferFilter, len(f.TransferFilter))
		for i := range f.TransferFilter {
			req.Transfers[i] = f.TransferFilter[i].ToGrpcFilter()
		}
	}
	if len(f.StorageDiffFilter) > 0 {
		req.StorageDiffs = make([]*pb.StorageDiffFilter, len(f.StorageDiffFilter))
		for i := range f.StorageDiffFilter {
			req.StorageDiffs[i] = f.StorageDiffFilter[i].ToGrpcFilter()
		}
	}
	if len(f.TokenBalanceFilter) > 0 {
		req.TokenBalances = make([]*pb.TokenBalanceFilter, len(f.TokenBalanceFilter))
		for i := range f.TokenBalanceFilter {
			req.TokenBalances[i] = f.TokenBalanceFilter[i].ToGrpcFilter()
		}
	}
	if len(f.TokenFilter) > 0 {
		req.Tokens = make([]*pb.TokenFilter, len(f.TokenFilter))
		for i := range f.TokenFilter {
			req.Tokens[i] = f.TokenFilter[i].ToGrpcFilter()
		}
	}

	return req
}

// AddressFilter -
type AddressFilter struct {
	Id           *IntegerFilter `yaml:"id" validate:"omitempty"`
	Height       *IntegerFilter `yaml:"height" validate:"omitempty"`
	OnlyStarknet bool           `yaml:"only_starknet" validate:"omitempty"`
}

// ToGrpcFilter -
func (f AddressFilter) ToGrpcFilter() *pb.AddressFilter {
	fltr := new(pb.AddressFilter)
	fltr.OnlyStarknet = f.OnlyStarknet

	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	return fltr
}

// DeclareFilters -
type DeclareFilters struct {
	Height  *IntegerFilter `yaml:"height" validate:"omitempty"`
	Time    *TimeFilter    `yaml:"time" validate:"omitempty"`
	Status  *EnumFilter    `yaml:"status" validate:"omitempty"`
	Version *EnumFilter    `yaml:"version" validate:"omitempty"`
	Id      *IntegerFilter `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f DeclareFilters) ToGrpcFilter() *pb.DeclareFilters {
	fltr := new(pb.DeclareFilters)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.Version != nil {
		fltr.Version = f.Version.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	return fltr
}

// DeployFilters -
type DeployFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Class          *BytesFilter      `yaml:"class" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f DeployFilters) ToGrpcFilter() *pb.DeployFilters {
	fltr := new(pb.DeployFilters)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Class != nil {
		fltr.Class = f.Class.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	return fltr
}

// DeployAccountFilters -
type DeployAccountFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Class          *BytesFilter      `yaml:"class" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f DeployAccountFilters) ToGrpcFilter() *pb.DeployAccountFilters {
	fltr := new(pb.DeployAccountFilters)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Class != nil {
		fltr.Class = f.Class.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	return fltr
}

// L1HandlerFilters -
type L1HandlerFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Contract       *BytesFilter      `yaml:"contract" validate:"omitempty"`
	Selector       *EqualityFilter   `yaml:"selector" validate:"omitempty"`
	Entrypoint     *StringFilter     `yaml:"entrypoint" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f L1HandlerFilters) ToGrpcFilter() *pb.L1HandlerFilter {
	fltr := new(pb.L1HandlerFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	if f.Entrypoint != nil {
		fltr.Entrypoint = f.Entrypoint.ToGrpcFilter()
	}
	if f.Selector != nil {
		fltr.Selector = f.Selector.ToGrpcFilter()
	}
	return fltr
}

// InvokeFilters -
type InvokeFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Version        *EnumFilter       `yaml:"version" validate:"omitempty"`
	Contract       *BytesFilter      `yaml:"contract" validate:"omitempty"`
	Selector       *EqualityFilter   `yaml:"selector" validate:"omitempty"`
	Entrypoint     *StringFilter     `yaml:"entrypoint" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f InvokeFilters) ToGrpcFilter() *pb.InvokeFilters {
	fltr := new(pb.InvokeFilters)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.Version != nil {
		fltr.Version = f.Version.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	if f.Entrypoint != nil {
		fltr.Entrypoint = f.Entrypoint.ToGrpcFilter()
	}
	if f.Selector != nil {
		fltr.Selector = f.Selector.ToGrpcFilter()
	}
	return fltr
}

// InternalFilters -
type InternalFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Contract       *BytesFilter      `yaml:"contract" validate:"omitempty"`
	Caller         *BytesFilter      `yaml:"caller" validate:"omitempty"`
	Class          *BytesFilter      `yaml:"class" validate:"omitempty"`
	Selector       *EqualityFilter   `yaml:"selector" validate:"omitempty"`
	Entrypoint     *StringFilter     `yaml:"entrypoint" validate:"omitempty"`
	EntrypointType *EnumFilter       `yaml:"entrypoint_type" validate:"omitempty"`
	CallType       *EnumFilter       `yaml:"call_type" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f InternalFilters) ToGrpcFilter() *pb.InternalFilter {
	fltr := new(pb.InternalFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Caller != nil {
		fltr.Caller = f.Caller.ToGrpcFilter()
	}
	if f.Class != nil {
		fltr.Class = f.Class.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.CallType != nil {
		fltr.CallType = f.CallType.ToGrpcFilter()
	}
	if f.EntrypointType != nil {
		fltr.CallType = f.EntrypointType.ToGrpcFilter()
	}
	if f.Entrypoint != nil {
		fltr.Entrypoint = f.Entrypoint.ToGrpcFilter()
	}
	if f.Selector != nil {
		fltr.Selector = f.Selector.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	return fltr
}

// FeeFilters -
type FeeFilters struct {
	Height         *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time           *TimeFilter       `yaml:"time" validate:"omitempty"`
	Status         *EnumFilter       `yaml:"status" validate:"omitempty"`
	Contract       *BytesFilter      `yaml:"contract" validate:"omitempty"`
	Caller         *BytesFilter      `yaml:"caller" validate:"omitempty"`
	Class          *BytesFilter      `yaml:"class" validate:"omitempty"`
	Selector       *EqualityFilter   `yaml:"selector" validate:"omitempty"`
	Entrypoint     *StringFilter     `yaml:"entrypoint" validate:"omitempty"`
	EntrypointType *EnumFilter       `yaml:"entrypoint_type" validate:"omitempty"`
	CallType       *EnumFilter       `yaml:"call_type" validate:"omitempty"`
	ParsedCalldata map[string]string `yaml:"parsed_calldata" validate:"omitempty"`
	Id             *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f FeeFilters) ToGrpcFilter() *pb.FeeFilter {
	fltr := new(pb.FeeFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Caller != nil {
		fltr.Caller = f.Caller.ToGrpcFilter()
	}
	if f.Class != nil {
		fltr.Class = f.Class.ToGrpcFilter()
	}
	if f.Status != nil {
		fltr.Status = f.Status.ToGrpcFilter()
	}
	if f.CallType != nil {
		fltr.CallType = f.CallType.ToGrpcFilter()
	}
	if f.EntrypointType != nil {
		fltr.CallType = f.EntrypointType.ToGrpcFilter()
	}
	if f.Entrypoint != nil {
		fltr.Entrypoint = f.Entrypoint.ToGrpcFilter()
	}
	if f.Selector != nil {
		fltr.Selector = f.Selector.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedCalldata != nil {
		fltr.ParsedCalldata = f.ParsedCalldata
	}
	return fltr
}

// EventFilter -
type EventFilter struct {
	Height     *IntegerFilter    `yaml:"height" validate:"omitempty"`
	Time       *TimeFilter       `yaml:"time" validate:"omitempty"`
	Contract   *BytesFilter      `yaml:"contract" validate:"omitempty"`
	From       *BytesFilter      `yaml:"from" validate:"omitempty"`
	Name       *StringFilter     `yaml:"name" validate:"omitempty"`
	ParsedData map[string]string `yaml:"parsed_data" validate:"omitempty"`
	Id         *IntegerFilter    `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f EventFilter) ToGrpcFilter() *pb.EventFilter {
	fltr := new(pb.EventFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.From != nil {
		fltr.From = f.From.ToGrpcFilter()
	}
	if f.Name != nil {
		fltr.Name = f.Name.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.ParsedData != nil {
		fltr.ParsedData = f.ParsedData
	}
	return fltr
}

// MessageFilter -
type MessageFilter struct {
	Height   *IntegerFilter  `yaml:"height" validate:"omitempty"`
	Time     *TimeFilter     `yaml:"time" validate:"omitempty"`
	Contract *BytesFilter    `yaml:"contract" validate:"omitempty"`
	From     *BytesFilter    `yaml:"from" validate:"omitempty"`
	To       *BytesFilter    `yaml:"to" validate:"omitempty"`
	Selector *EqualityFilter `yaml:"selector" validate:"omitempty"`
	Id       *IntegerFilter  `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f MessageFilter) ToGrpcFilter() *pb.MessageFilter {
	fltr := new(pb.MessageFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.From != nil {
		fltr.From = f.From.ToGrpcFilter()
	}
	if f.To != nil {
		fltr.To = f.To.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.Selector != nil {
		fltr.Selector = f.Selector.ToGrpcFilter()
	}
	return fltr
}

// TransferFilter -
type TransferFilter struct {
	Height   *IntegerFilter `yaml:"height" validate:"omitempty"`
	Time     *TimeFilter    `yaml:"time" validate:"omitempty"`
	Contract *BytesFilter   `yaml:"contract" validate:"omitempty"`
	From     *BytesFilter   `yaml:"from" validate:"omitempty"`
	To       *BytesFilter   `yaml:"to" validate:"omitempty"`
	TokenId  *StringFilter  `yaml:"token_id" validate:"omitempty"`
	Id       *IntegerFilter `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f TransferFilter) ToGrpcFilter() *pb.TransferFilter {
	fltr := new(pb.TransferFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Time != nil {
		fltr.Time = f.Time.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.From != nil {
		fltr.From = f.From.ToGrpcFilter()
	}
	if f.To != nil {
		fltr.To = f.To.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.TokenId != nil {
		fltr.TokenId = f.TokenId.ToGrpcFilter()
	}
	return fltr
}

// StorageDiffFilter -
type StorageDiffFilter struct {
	Height   *IntegerFilter  `yaml:"height" validate:"omitempty"`
	Contract *BytesFilter    `yaml:"contract" validate:"omitempty"`
	Key      *EqualityFilter `yaml:"key" validate:"omitempty"`
	Id       *IntegerFilter  `yaml:"id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f StorageDiffFilter) ToGrpcFilter() *pb.StorageDiffFilter {
	fltr := new(pb.StorageDiffFilter)

	if f.Height != nil {
		fltr.Height = f.Height.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Key != nil {
		fltr.Key = f.Key.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	return fltr
}

// TokenBalanceFilter -
type TokenBalanceFilter struct {
	Owner    *BytesFilter  `yaml:"owner" validate:"omitempty"`
	Contract *BytesFilter  `yaml:"contract" validate:"omitempty"`
	TokenId  *StringFilter `yaml:"token_id" validate:"omitempty"`
}

// ToGrpcFilter -
func (f TokenBalanceFilter) ToGrpcFilter() *pb.TokenBalanceFilter {
	fltr := new(pb.TokenBalanceFilter)

	if f.Owner != nil {
		fltr.Owner = f.Owner.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.TokenId != nil {
		fltr.TokenId = f.TokenId.ToGrpcFilter()
	}
	return fltr
}

// TokenFilter -
type TokenFilter struct {
	TokenId  *StringFilter     `yaml:"token_id" validate:"omitempty"`
	Contract *BytesFilter      `yaml:"contract" validate:"omitempty"`
	Id       *IntegerFilter    `yaml:"id" validate:"omitempty"`
	Type     *EnumStringFilter `yaml:"type" validate:"omitempty"`
}

// ToGrpcFilter -
func (f TokenFilter) ToGrpcFilter() *pb.TokenFilter {
	fltr := new(pb.TokenFilter)

	if f.TokenId != nil {
		fltr.TokenId = f.TokenId.ToGrpcFilter()
	}
	if f.Contract != nil {
		fltr.Contract = f.Contract.ToGrpcFilter()
	}
	if f.Id != nil {
		fltr.Id = f.Id.ToGrpcFilter()
	}
	if f.Type != nil {
		fltr.Type = f.Type.ToGrpcFilter()
	}
	return fltr
}

// IntegerFilter -
type IntegerFilter struct {
	Gt      uint64          `yaml:"gt" validate:"omitempty"`
	Gte     uint64          `yaml:"gte" validate:"omitempty"`
	Lt      uint64          `yaml:"lt" validate:"omitempty"`
	Lte     uint64          `yaml:"lte" validate:"omitempty"`
	Eq      uint64          `yaml:"eq" validate:"omitempty"`
	Neq     uint64          `yaml:"neq" validate:"omitempty"`
	Between *BetweenInteger `yaml:"between" validate:"omitempty"`
}

// ToGrpcFilter -
func (f IntegerFilter) ToGrpcFilter() *pb.IntegerFilter {
	fltr := new(pb.IntegerFilter)

	switch {
	case f.Between != nil:
		fltr.Filter = &pb.IntegerFilter_Between{
			Between: &pb.BetweenInteger{
				From: f.Between.From,
				To:   f.Between.To,
			},
		}
	case f.Gt > 0:
		fltr.Filter = &pb.IntegerFilter_Gt{
			Gt: f.Gt,
		}
	case f.Gte > 0:
		fltr.Filter = &pb.IntegerFilter_Gte{
			Gte: f.Gte,
		}
	case f.Lte > 0:
		fltr.Filter = &pb.IntegerFilter_Lte{
			Lte: f.Lte,
		}
	case f.Lt > 0:
		fltr.Filter = &pb.IntegerFilter_Lt{
			Lt: f.Lt,
		}
	case f.Eq > 0:
		fltr.Filter = &pb.IntegerFilter_Eq{
			Eq: f.Eq,
		}
	case f.Neq > 0:
		fltr.Filter = &pb.IntegerFilter_Neq{
			Neq: f.Neq,
		}
	}

	return fltr
}

// TimeFilter -
type TimeFilter struct {
	Gt      uint64          `yaml:"gt" validate:"omitempty"`
	Gte     uint64          `yaml:"gte" validate:"omitempty"`
	Lt      uint64          `yaml:"lt" validate:"omitempty"`
	Lte     uint64          `yaml:"lte" validate:"omitempty"`
	Between *BetweenInteger `yaml:"between" validate:"omitempty"`
}

// ToGrpcFilter -
func (f TimeFilter) ToGrpcFilter() *pb.TimeFilter {
	fltr := new(pb.TimeFilter)

	switch {
	case f.Between != nil:
		fltr.Filter = &pb.TimeFilter_Between{
			Between: &pb.BetweenInteger{
				From: f.Between.From,
				To:   f.Between.To,
			},
		}
	case f.Gt > 0:
		fltr.Filter = &pb.TimeFilter_Gt{
			Gt: f.Gt,
		}
	case f.Gte > 0:
		fltr.Filter = &pb.TimeFilter_Gte{
			Gte: f.Gte,
		}
	case f.Lte > 0:
		fltr.Filter = &pb.TimeFilter_Lte{
			Lte: f.Lte,
		}
	case f.Lt > 0:
		fltr.Filter = &pb.TimeFilter_Lt{
			Lt: f.Lt,
		}

	}

	return fltr
}

// Bytes -
type Bytes []byte

// UnmarshalYAML -
func (b *Bytes) UnmarshalYAML(node *yaml.Node) error {
	value := strings.TrimPrefix(node.Value, "0x")
	ba, err := hex.DecodeString(value)
	if err != nil {
		return err
	}
	*b = ba
	return nil
}

// BytesFilter -
type BytesFilter struct {
	Eq Bytes   `yaml:"eq" validate:"omitempty"`
	In []Bytes `yaml:"in" validate:"omitempty"`
}

// ToGrpcFilter -
func (f BytesFilter) ToGrpcFilter() *pb.BytesFilter {
	fltr := new(pb.BytesFilter)
	switch {
	case len(f.Eq) > 0:
		fltr.Filter = &pb.BytesFilter_Eq{
			Eq: f.Eq,
		}
	case len(f.In) > 0:
		in := make([][]byte, len(f.In))
		for i := range f.In {
			in[i] = []byte(f.In[i])
		}
		fltr.Filter = &pb.BytesFilter_In{
			In: &pb.BytesArray{
				Arr: in,
			},
		}
	}
	return fltr
}

// StringFilter -
type StringFilter struct {
	Eq string   `yaml:"eq" validate:"omitempty"`
	In []string `yaml:"in" validate:"omitempty"`
}

// ToGrpcFilter -
func (f StringFilter) ToGrpcFilter() *pb.StringFilter {
	fltr := new(pb.StringFilter)

	switch {
	case f.Eq != "":
		fltr.Filter = &pb.StringFilter_Eq{
			Eq: f.Eq,
		}
	case len(f.In) > 0:
		fltr.Filter = &pb.StringFilter_In{
			In: &pb.StringArray{
				Arr: f.In,
			},
		}
	}
	return fltr
}

// BetweenInteger -
type BetweenInteger struct {
	From uint64 `yaml:"from" validate:"required"`
	To   uint64 `yaml:"to" validate:"required"`
}

// ToGrpcFilter -
func (f BetweenInteger) ToGrpcFilter() *pb.BetweenInteger {
	fltr := new(pb.BetweenInteger)
	fltr.From = f.From
	fltr.To = f.To
	return fltr
}

// EnumFilter -
type EnumFilter struct {
	Eq    uint64   `yaml:"eq" validate:"omitempty"`
	Neq   uint64   `yaml:"neq" validate:"omitempty"`
	In    []uint64 `yaml:"in" validate:"omitempty"`
	Notin []uint64 `yaml:"notin" validate:"omitempty"`
}

// ToGrpcFilter -
func (f EnumFilter) ToGrpcFilter() *pb.EnumFilter {
	fltr := new(pb.EnumFilter)
	switch {
	case f.Eq > 0:
		fltr.Filter = &pb.EnumFilter_Eq{
			Eq: f.Eq,
		}
	case f.Neq > 0:
		fltr.Filter = &pb.EnumFilter_Neq{
			Neq: f.Neq,
		}
	case len(f.In) > 0:
		fltr.Filter = &pb.EnumFilter_In{
			In: &pb.IntegerArray{
				Arr: f.In,
			},
		}
	case len(f.Notin) > 0:
		fltr.Filter = &pb.EnumFilter_Notin{
			Notin: &pb.IntegerArray{
				Arr: f.Notin,
			},
		}
	}
	return fltr
}

// EnumStringFilter -
type EnumStringFilter struct {
	Eq    string   `yaml:"eq" validate:"omitempty"`
	Neq   string   `yaml:"neq" validate:"omitempty"`
	In    []string `yaml:"in" validate:"omitempty"`
	Notin []string `yaml:"notin" validate:"omitempty"`
}

// ToGrpcFilter -
func (f EnumStringFilter) ToGrpcFilter() *pb.EnumStringFilter {
	fltr := new(pb.EnumStringFilter)
	switch {
	case f.Eq != "":
		fltr.Filter = &pb.EnumStringFilter_Eq{
			Eq: f.Eq,
		}
	case f.Neq != "":
		fltr.Filter = &pb.EnumStringFilter_Neq{
			Neq: f.Neq,
		}
	case len(f.In) > 0:
		fltr.Filter = &pb.EnumStringFilter_In{
			In: &pb.StringArray{
				Arr: f.In,
			},
		}
	case len(f.Notin) > 0:
		fltr.Filter = &pb.EnumStringFilter_Notin{
			Notin: &pb.StringArray{
				Arr: f.Notin,
			},
		}
	}
	return fltr
}

// EqualityFilter -
type EqualityFilter struct {
	Eq  string `yaml:"eq" validate:"omitempty"`
	Neq string `yaml:"neq" validate:"omitempty"`
}

// ToGrpcFilter -
func (f EqualityFilter) ToGrpcFilter() *pb.EqualityFilter {
	fltr := new(pb.EqualityFilter)
	switch {
	case f.Eq != "":
		fltr.Filter = &pb.EqualityFilter_Eq{
			Eq: f.Eq,
		}
	case f.Neq != "":
		fltr.Filter = &pb.EqualityFilter_Neq{
			Neq: f.Neq,
		}
	}
	return fltr
}
