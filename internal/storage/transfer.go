package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// ITransfer -
type ITransfer interface {
	storage.Table[*Transfer]

	Copiable[Transfer]
	Filterable[Transfer, TransferFilter]
}

// TransferFilter -
type TransferFilter struct {
	ID       IntegerFilter
	Height   IntegerFilter
	Time     TimeFilter
	Contract BytesFilter
	From     BytesFilter
	To       BytesFilter
	TokenId  StringFilter
}

// Transfer -
type Transfer struct {
	// nolint
	tableName struct{} `pg:"transfer,partition_by:RANGE(time)"`

	ID         uint64          `pg:",pk,comment:Unique internal identity"`
	Height     uint64          `pg:",use_zero,comment:Block height"`
	Time       time.Time       `pg:",pk,comment:Block time"`
	ContractID uint64          `pg:",comment:Token contract id"`
	FromID     uint64          `pg:",comment:Id address which transfer from"`
	ToID       uint64          `pg:",comment:Id address which transfer to"`
	Amount     decimal.Decimal `pg:",type:numeric,use_zero,comment:Amount of transfer"`
	TokenID    decimal.Decimal `pg:",type:numeric,use_zero,comment:Id token which was transferred"`

	InvokeID        *uint64 `pg:",comment:Parent invoke id"`
	DeclareID       *uint64 `pg:",comment:Parent declare id"`
	DeployID        *uint64 `pg:",comment:Parent deploy id"`
	DeployAccountID *uint64 `pg:",comment:Parent deploy account id"`
	L1HandlerID     *uint64 `pg:",comment:Parent l1 handler id"`
	FeeID           *uint64 `pg:",comment:Parent fee invocation id"`
	InternalID      *uint64 `pg:",comment:Parent internal transaction id"`

	From     Address `pg:"rel:has-one" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	To       Address `pg:"rel:has-one" hasura:"table:address,field:to_id,remote_field:id,type:oto,name:to"`
	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Token    Token   `pg:"rel:has-one,fk:contract_id" hasura:"table:token,field:contract_id,remote_field:contract_id,type:oto,name:token"`
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

// GetHeight -
func (transfer Transfer) GetHeight() uint64 {
	return transfer.Height
}
