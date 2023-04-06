package storage

import (
	"github.com/dipdup-io/starknet-go-api/pkg/data"
)

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
	case data.TransactionTypeInvoke:
		return TransactionTypeInvoke
	case data.TransactionTypeDeclare:
		return TransactionTypeDeclare
	case data.TransactionTypeDeploy:
		return TransactionTypeDeploy
	case data.TransactionTypeDeployAccount:
		return TransactionTypeDeployAccount
	case data.TransactionTypeL1Handler:
		return TransactionTypeL1Handler
	default:
		return TransactionTypeUnknown
	}
}

// String -
func (t TransactionType) String() string {
	switch t {
	case TransactionTypeDeclare:
		return data.TransactionTypeDeclare
	case TransactionTypeDeploy:
		return data.TransactionTypeDeploy
	case TransactionTypeDeployAccount:
		return data.TransactionTypeDeployAccount
	case TransactionTypeInvoke:
		return data.TransactionTypeInvoke
	case TransactionTypeL1Handler:
		return data.TransactionTypeL1Handler
	default:
		return Unknown
	}
}
