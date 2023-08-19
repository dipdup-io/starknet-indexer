package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// ITransfer -
type ITransfer interface {
	storage.Table[*Transfer]
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
	bun.BaseModel `bun:"transfer" comment:"Trandfer table" partition:"RANGE(time)"`

	ID         uint64          `bun:",pk,autoincrement" comment:"Unique internal identity"`
	Height     uint64          `comment:"Block height"`
	Time       time.Time       `bun:",pk" comment:"Block time"`
	ContractID uint64          `comment:"Token contract id"`
	FromID     uint64          `comment:"Id address which transfer from"`
	ToID       uint64          `comment:"Id address which transfer to"`
	Amount     decimal.Decimal `bun:",type:numeric" comment:"Amount of transfer"`
	TokenID    decimal.Decimal `bun:",type:numeric" comment:"Id token which was transferred"`

	InvokeID        *uint64 `comment:"Parent invoke id"`
	DeclareID       *uint64 `comment:"Parent declare id"`
	DeployID        *uint64 `comment:"Parent deploy id"`
	DeployAccountID *uint64 `comment:"Parent deploy account id"`
	L1HandlerID     *uint64 `comment:"Parent l1 handler id"`
	FeeID           *uint64 `comment:"Parent fee invocation id"`
	InternalID      *uint64 `comment:"Parent internal transaction id"`

	From     Address `bun:"rel:belongs-to" hasura:"table:address,field:from_id,remote_field:id,type:oto,name:from"`
	To       Address `bun:"rel:belongs-to" hasura:"table:address,field:to_id,remote_field:id,type:oto,name:to"`
	Contract Address `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Token    Token   `bun:"-" hasura:"table:token,field:contract_id,remote_field:contract_id,type:oto,name:token"`
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

// GetId -
func (transfer Transfer) GetId() uint64 {
	return transfer.ID
}

// Columns -
func (Transfer) Columns() []string {
	return []string{
		"height", "time", "contract_id", "from_id",
		"to_id", "amount", "token_id", "invoke_id",
		"declare_id", "deploy_id", "deploy_account_id",
		"l1_handler_id", "fee_id", "internal_id",
	}
}

// Flat -
func (transfer Transfer) Flat() []any {
	return []any{
		transfer.Height,
		transfer.Time,
		transfer.ContractID,
		transfer.FromID,
		transfer.ToID,
		transfer.Amount,
		transfer.TokenID,
		transfer.InvokeID,
		transfer.DeclareID,
		transfer.DeployID,
		transfer.DeployAccountID,
		transfer.L1HandlerID,
		transfer.FeeID,
		transfer.InternalID,
	}
}
