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
	Subscriptions map[string]Subscription `yaml:"subscriptions" validate:"required"`
}

// TODO: implement all filters
// Subscription -
type Subscription struct {
	Head                 bool                     `yaml:"head" validate:"omitempty"`
	InvokeFilters        *pb.InvokeFilters        `yaml:"invokes" validate:"omitempty"`
	DeclareFilters       *pb.DeclareFilters       `yaml:"filters" validate:"omitempty"`
	DeployFilters        *pb.DeployFilters        `yaml:"deploys" validate:"omitempty"`
	DeployAccountFilters *pb.DeployAccountFilters `yaml:"deploy_accounts" validate:"omitempty"`
	L1HandlerFilter      *pb.L1HandlerFilter      `yaml:"l1_handlers" validate:"omitempty"`
	InternalFilter       *pb.InternalFilter       `yaml:"internals" validate:"omitempty"`
	FeeFilter            *pb.FeeFilter            `yaml:"fees" validate:"omitempty"`
	EventFilter          *EventFilter             `yaml:"events" validate:"omitempty"`
	MessageFilter        *pb.MessageFilter        `yaml:"messages" validate:"omitempty"`
	TransferFilter       *pb.TransferFilter       `yaml:"transfers" validate:"omitempty"`
	StorageDiffFilter    *pb.StorageDiffFilter    `yaml:"storage_diffs" validate:"omitempty"`
	TokenBalanceFilter   *pb.TokenBalanceFilter   `yaml:"token_balances" validate:"omitempty"`
}

// ToGrpcFilter -
func (f Subscription) ToGrpcFilter() *pb.SubscribeRequest {
	req := new(pb.SubscribeRequest)
	req.Head = f.Head

	if f.EventFilter != nil {
		req.Events = f.EventFilter.ToGrpcFilter()
	}

	return req
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
