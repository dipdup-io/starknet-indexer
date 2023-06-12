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
	FeeID           *uint64
	InternalID      *uint64

	Invoke        *storage.Invoke
	Declare       *storage.Declare
	Deploy        *storage.Deploy
	DeployAccount *storage.DeployAccount
	L1Handler     *storage.L1Handler
	Fee           *storage.Fee
	Internal      *storage.Internal

	Height     uint64
	Time       time.Time
	Status     storage.Status
	Hash       []byte
	ProxyId    uint64
	ContractId uint64

	ProxyUpgrades ProxyMap[struct{}]
}

// NewEmptyTxContext -
func NewEmptyTxContext() TxContext {
	return TxContext{
		ProxyUpgrades: NewProxyMap[struct{}](),
	}
}

// NewTxContextFromInvoke -
func NewTxContextFromInvoke(tx storage.Invoke, proxyId uint64) TxContext {
	return TxContext{
		InvokeID:      &tx.ID,
		Invoke:        &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		Hash:          tx.Hash,
		ProxyId:       proxyId,
		ContractId:    tx.ContractID,
		ProxyUpgrades: NewProxyMap[struct{}](),
	}
}

// NewTxContextFromDeclare -
func NewTxContextFromDeclare(tx storage.Declare, proxyId uint64) TxContext {
	var contractId uint64
	if tx.ContractID != nil {
		contractId = *tx.ContractID
	}
	return TxContext{
		DeclareID:     &tx.ID,
		Declare:       &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		Hash:          tx.Hash,
		ProxyId:       proxyId,
		ContractId:    contractId,
		ProxyUpgrades: NewProxyMap[struct{}](),
	}
}

// NewTxContextFromDeploy -
func NewTxContextFromDeploy(tx storage.Deploy, proxyId uint64) TxContext {
	return TxContext{
		DeployID:      &tx.ID,
		Deploy:        &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		Hash:          tx.Hash,
		ProxyId:       proxyId,
		ContractId:    tx.ContractID,
		ProxyUpgrades: NewProxyMap[struct{}](),
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
		ProxyUpgrades:   NewProxyMap[struct{}](),
	}
}

// NewTxContextFromL1Hadler -
func NewTxContextFromL1Hadler(tx storage.L1Handler, proxyId uint64) TxContext {
	return TxContext{
		L1HandlerID:   &tx.ID,
		L1Handler:     &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		Hash:          tx.Hash,
		ProxyId:       proxyId,
		ContractId:    tx.ContractID,
		ProxyUpgrades: NewProxyMap[struct{}](),
	}
}

// NewTxContextFromInternal -
func NewTxContextFromInternal(tx storage.Internal, proxyUpgrades ProxyMap[struct{}], proxyId uint64) TxContext {
	cloneProxyMap, _ := proxyUpgrades.Clone()
	return TxContext{
		InternalID:    &tx.ID,
		Internal:      &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		Hash:          tx.Hash,
		ProxyId:       proxyId,
		ContractId:    tx.ContractID,
		ProxyUpgrades: cloneProxyMap,
	}
}

// NewTxContextFromFee -
func NewTxContextFromFee(tx storage.Fee, proxyId uint64) TxContext {
	return TxContext{
		FeeID:         &tx.ID,
		Fee:           &tx,
		Height:        tx.Height,
		Time:          tx.Time,
		Status:        tx.Status,
		ProxyId:       proxyId,
		ContractId:    tx.ContractID,
		ProxyUpgrades: NewProxyMap[struct{}](),
	}
}
