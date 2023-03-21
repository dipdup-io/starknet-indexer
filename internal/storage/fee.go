package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IFee -
type IFee interface {
	storage.Table[*Fee]
}

// Fee -
type Fee struct {
	// nolint
	tableName struct{} `pg:"fee,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`
	Status Status    `pg:",use_zero"`

	ContractID uint64
	FromID     uint64
	ToID       uint64
	Amount     decimal.Decimal `pg:",type:numeric,use_zero"`

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeID        *uint64
	L1HandlerID     *uint64

	From     Address `pg:"rel:has-one"`
	To       Address `pg:"rel:has-one"`
	Contract Address `pg:"rel:has-one"`
}

// TableName -
func (Fee) TableName() string {
	return "fee"
}

// TokenBalanceUpdates -
func (fee Fee) TokenBalanceUpdates() []TokenBalance {
	updates := make([]TokenBalance, 0)
	if fee.FromID > 0 {
		updates = append(updates, TokenBalance{
			OwnerID:    fee.FromID,
			ContractID: fee.ContractID,
			TokenID:    decimal.Zero,
			Balance:    fee.Amount.Copy().Neg(),
		})
	}

	if fee.ToID > 0 {
		updates = append(updates, TokenBalance{
			OwnerID:    fee.ToID,
			ContractID: fee.ContractID,
			TokenID:    decimal.Zero,
			Balance:    fee.Amount.Copy(),
		})
	}

	return updates
}
