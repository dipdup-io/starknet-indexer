package storage

import (
	"context"
	"fmt"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// ITokenBalance -
type ITokenBalance interface {
	storage.Table[*TokenBalance]
	Filterable[TokenBalance, TokenBalanceFilter]

	Balances(ctx context.Context, contractId uint64, tokenId int64, limit, offset int) ([]TokenBalance, error)
}

// TokenBalanceFilter -
type TokenBalanceFilter struct {
	Owner    BytesFilter
	Contract BytesFilter
	TokenId  StringFilter
}

// TokenBalance -
type TokenBalance struct {
	bun.BaseModel `bun:"token_balance" comment:"Table with token balances"`

	OwnerID    uint64          `bun:",pk" comment:"Identity of owner address"`
	ContractID uint64          `bun:",pk" comment:"Identity of contract address"`
	TokenID    decimal.Decimal `bun:",pk,type:numeric" comment:"Token id"`
	Balance    decimal.Decimal `bun:",type:numeric" comment:"Token balance"`

	Contract Address `bun:"rel:belongs-to,join:contract_id=id" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `bun:"rel:belongs-to,join:owner_id=id" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
}

// TableName -
func (TokenBalance) TableName() string {
	return "token_balance"
}

// String -
func (tb TokenBalance) String() string {
	return fmt.Sprintf("%d_%d_%s", tb.ContractID, tb.OwnerID, tb.TokenID.String())
}
