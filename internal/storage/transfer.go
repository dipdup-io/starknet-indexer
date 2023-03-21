package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// ITransfer -
type ITransfer interface {
	storage.Table[*Transfer]
}

// Transfer -
type Transfer struct {
	// nolint
	tableName struct{} `pg:"transfer,partition_by:RANGE(time)"`

	ID         uint64    `pg:",pk"`
	Height     uint64    `pg:",use_zero"`
	Time       time.Time `pg:",pk"`
	ContractID uint64
	FromID     uint64
	ToID       uint64
	Amount     decimal.Decimal `pg:",type:numeric,use_zero"`
	TokenID    decimal.Decimal `pg:",type:numeric,use_zero"`

	InvokeID        *uint64
	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	L1HandlerID     *uint64
	InternalID      *uint64

	From     Address `pg:"rel:has-one"`
	To       Address `pg:"rel:has-one"`
	Contract Address `pg:"rel:has-one"`
}

// TableName -
func (Transfer) TableName() string {
	return "transfer"
}

// TokenBalanceUpdates -
func (transfer Transfer) TokenBalanceUpdates() []TokenBalance {
	if transfer.FromID == transfer.ToID {
		return nil
	}
	updates := make([]TokenBalance, 0)
	if transfer.FromID > 0 {
		updates = append(updates, TokenBalance{
			OwnerID:    transfer.FromID,
			ContractID: transfer.ContractID,
			TokenID:    transfer.TokenID,
			Balance:    transfer.Amount.Copy().Neg(),
		})
	}

	if transfer.ToID > 0 {
		updates = append(updates, TokenBalance{
			OwnerID:    transfer.ToID,
			ContractID: transfer.ContractID,
			TokenID:    transfer.TokenID,
			Balance:    transfer.Amount.Copy(),
		})
	}

	return updates
}
