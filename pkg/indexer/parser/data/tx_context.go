package data

import (
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
)

// TxContext -
type TxContext struct {
	InvokeID        *uint64
	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	L1HandlerID     *uint64
	InternalID      *uint64

	Invoke        *storage.Invoke
	Declare       *storage.Declare
	Deploy        *storage.Deploy
	DeployAccount *storage.DeployAccount
	L1Handler     *storage.L1Handler
	Internal      *storage.Internal

	Height     uint64
	Time       time.Time
	Status     storage.Status
	Hash       []byte
	ProxyId    uint64
	ContractId uint64

	TransfersCount int
}

// NewTxContextFromInvoke -
func NewTxContextFromInvoke(tx storage.Invoke, proxyId uint64) TxContext {
	return TxContext{
		InvokeID:       &tx.ID,
		Invoke:         &tx,
		Height:         tx.Height,
		Time:           tx.Time,
		Status:         tx.Status,
		Hash:           tx.Hash,
		ProxyId:        proxyId,
		ContractId:     tx.ContractID,
		TransfersCount: len(tx.Transfers),
	}
}

// NewTxContextFromDeclare -
func NewTxContextFromDeclare(tx storage.Declare, proxyId uint64) TxContext {
	var contractId uint64
	if tx.ContractID != nil {
		contractId = *tx.ContractID
	}
	return TxContext{
		DeclareID:  &tx.ID,
		Declare:    &tx,
		Height:     tx.Height,
		Time:       tx.Time,
		Status:     tx.Status,
		Hash:       tx.Hash,
		ProxyId:    proxyId,
		ContractId: contractId,
	}
}

// NewTxContextFromDeploy -
func NewTxContextFromDeploy(tx storage.Deploy, proxyId uint64) TxContext {
	return TxContext{
		DeployID:   &tx.ID,
		Deploy:     &tx,
		Height:     tx.Height,
		Time:       tx.Time,
		Status:     tx.Status,
		Hash:       tx.Hash,
		ProxyId:    proxyId,
		ContractId: tx.ContractID,
	}
}

// NewTxContextFromDeployAccount -
func NewTxContextFromDeployAccount(tx storage.DeployAccount, proxyId uint64) TxContext {
	return TxContext{
		DeployAccountID: &tx.ID,
		DeployAccount:   &tx,
		Height:          tx.Height,
		Time:            tx.Time,
		Status:          tx.Status,
		Hash:            tx.Hash,
		ProxyId:         proxyId,
		ContractId:      tx.ContractID,
	}
}

// NewTxContextFromL1Hadler -
func NewTxContextFromL1Hadler(tx storage.L1Handler, proxyId uint64) TxContext {
	return TxContext{
		L1HandlerID: &tx.ID,
		L1Handler:   &tx,
		Height:      tx.Height,
		Time:        tx.Time,
		Status:      tx.Status,
		Hash:        tx.Hash,
		ProxyId:     proxyId,
		ContractId:  tx.ContractID,
	}
}

// NewTxContextFromInternal -
func NewTxContextFromInternal(tx storage.Internal, proxyId uint64) TxContext {
	return TxContext{
		InternalID:     &tx.ID,
		Internal:       &tx,
		Height:         tx.Height,
		Time:           tx.Time,
		Status:         tx.Status,
		Hash:           tx.Hash,
		ProxyId:        proxyId,
		ContractId:     tx.ContractID,
		TransfersCount: len(tx.Transfers),
	}
}
