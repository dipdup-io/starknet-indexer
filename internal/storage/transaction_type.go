package storage

import starknet "github.com/dipdup-io/starknet-go-api/pkg/api"

// TransactionType -
type TransactionType int

const (
	TransactionTypeUnknown TransactionType = iota + 1
	TransactionTypeInvoke
	TransactionTypeDeclare
	TransactionTypeDeploy
	TransactionTypeDeployAccount
	TransactionTypeL1Handler
)

// NewTransactionType -
func NewTransactionType(value string) TransactionType {
	switch value {
	case starknet.TransactionTypeInvoke:
		return TransactionTypeInvoke
	case starknet.TransactionTypeDeclare:
		return TransactionTypeDeclare
	case starknet.TransactionTypeDeploy:
		return TransactionTypeDeploy
	case starknet.TransactionTypeDeployAccount:
		return TransactionTypeDeployAccount
	case starknet.TransactionTypeL1Handler:
		return TransactionTypeL1Handler
	default:
		return TransactionTypeUnknown
	}
}

// String -
func (t TransactionType) String() string {
	switch t {
	case TransactionTypeDeclare:
		return starknet.TransactionTypeDeclare
	case TransactionTypeDeploy:
		return starknet.TransactionTypeDeploy
	case TransactionTypeDeployAccount:
		return starknet.TransactionTypeDeployAccount
	case TransactionTypeInvoke:
		return starknet.TransactionTypeInvoke
	case TransactionTypeL1Handler:
		return starknet.TransactionTypeL1Handler
	default:
		return "UNKNOWN"
	}
}
