package decode

import (
	"testing"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	type args struct {
		contractAbi abi.Abi
		keys        []string
		data        []string
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]any
		want1 string
	}{
		{
			name: "transfer before v0.13.3",
			args: args{
				contractAbi: abi.Abi{
					EventsBySelector: map[string]*abi.EventItem{
						"99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9": {
							Type: abi.Type{
								Name: "Transfer",
								Type: "event",
							},
							Data: []abi.Type{
								{
									Name: "from",
									Type: "core::starknet::contract_address::ContractAddress",
								},
								{
									Name: "to",
									Type: "core::starknet::contract_address::ContractAddress",
								},
								{
									Name: "amount",
									Type: "core::integer::u256",
								},
							},
							Keys: []abi.Type{},
						},
					},
				},
				keys: []string{
					"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9",
				},
				data: []string{
					"0x3a39bfdf7aad9504978afc793d5c1e2c8d9fc6f2e02720aebc52c9817a83e42",
					"0x0",
					"0x21dcf820ad93f",
					"0x0",
				},
			},
			want: map[string]any{
				"from":   "0x3a39bfdf7aad9504978afc793d5c1e2c8d9fc6f2e02720aebc52c9817a83e42",
				"to":     "0x0",
				"amount": "595727030606143",
			},
			want1: "Transfer",
		}, {
			name: "transfer after v0.13.3",
			args: args{
				contractAbi: abi.Abi{
					EventsBySelector: map[string]*abi.EventItem{
						"99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9": {
							Type: abi.Type{
								Name: "Transfer",
								Type: "event",
							},
							Data: []abi.Type{
								{
									Name: "amount",
									Type: "core::integer::u256",
								},
							},
							Keys: []abi.Type{
								{
									Name: "from",
									Type: "core::starknet::contract_address::ContractAddress",
								},
								{
									Name: "to",
									Type: "core::starknet::contract_address::ContractAddress",
								},
							},
						},
					},
				},
				keys: []string{
					"0x99cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9",
					"0x3a39bfdf7aad9504978afc793d5c1e2c8d9fc6f2e02720aebc52c9817a83e42",
					"0x0",
				},
				data: []string{
					"0x21dcf820ad93f",
					"0x0",
				},
			},
			want: map[string]any{
				"from":   "0x3a39bfdf7aad9504978afc793d5c1e2c8d9fc6f2e02720aebc52c9817a83e42",
				"to":     "0x0",
				"amount": "595727030606143",
			},
			want1: "Transfer",
		}, {
			name: "test 3",
			args: args{
				contractAbi: abi.Abi{
					EventsBySelector: map[string]*abi.EventItem{
						"31f8daa2ac8dacd06ab968bad8f97f98f63c30a86dbfcebdd7625f589d4e7e6": {
							Type: abi.Type{
								Name: "hello::HelloStarknet::BalanceIncreased",
								Type: "event",
								Kind: "struct",
							},
							Members: []abi.Type{
								{
									Name: "amount",
									Type: "core::felt252",
									Kind: "data",
								}, {
									Name: "by",
									Type: "core::starknet::contract_address::ContractAddress",
									Kind: "key",
								},
							},
						},
					},
				},
				keys: []string{
					"0x31f8daa2ac8dacd06ab968bad8f97f98f63c30a86dbfcebdd7625f589d4e7e6",
					"0x1a62446e05ee60540d94b2e731ed037a1798065f9b8e719e293180b493b91f7",
				},
				data: []string{
					"0x2a",
				},
			},
			want: map[string]any{
				"by":     "0x1a62446e05ee60540d94b2e731ed037a1798065f9b8e719e293180b493b91f7",
				"amount": "0x2a",
			},
			want1: "hello::HelloStarknet::BalanceIncreased",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Event(tt.args.contractAbi, tt.args.keys, tt.args.data)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.want1, got1)
		})
	}
}
