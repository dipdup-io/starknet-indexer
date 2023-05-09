package storage

import (
	"context"
	"fmt"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// ITokenBalance -
type ITokenBalance interface {
	storage.Table[*TokenBalance]
	Filterable[TokenBalance, TokenBalanceFilter]

	NegativeBalances(ctx context.Context) ([]TokenBalance, error)
	TotalSupply(ctx context.Context, contractId, tokenId uint64) (decimal.Decimal, error)
	Owner(ctx context.Context, cotractId uint64, tokenId decimal.Decimal) (TokenBalance, error)
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
	// nolint
	tableName struct{} `pg:"token_balance,comment:Table with token balances"`

	OwnerID    uint64          `pg:",pk,comment:Identity of owner address"`
	ContractID uint64          `pg:",pk,comment:Identity of contract address"`
	TokenID    decimal.Decimal `pg:",pk,type:numeric,use_zero,comment:Token id"`
	Balance    decimal.Decimal `pg:",type:numeric,use_zero,comment:Token balance"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `pg:"rel:has-one" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
	Token    Token   `pg:"rel:has-one,fk:contract_id" hasura:"table:token,field:contract_id,remote_field:contract_id,type:oto,name:token"`
}

// TableName -
func (TokenBalance) TableName() string {
	return "token_balance"
}

// String -
func (tb TokenBalance) String() string {
	return fmt.Sprintf("%d_%d_%s", tb.ContractID, tb.OwnerID, tb.TokenID.String())
}
