package storage

import (
	"context"
	"fmt"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// TokenType -
type TokenType string

// token types
const (
	TokenTypeERC20   TokenType = "erc20"
	TokenTypeERC721  TokenType = "erc721"
	TokenTypeERC1155 TokenType = "erc1155"
)

// TokenFilter -
type TokenFilter struct {
	ID       IntegerFilter
	Contract IdFilter
	TokenId  StringFilter
	Type     EnumStringFilter
}

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IToken interface {
	storage.Table[*Token]
	Filterable[Token, TokenFilter]

	Find(ctx context.Context, contractId uint64, tokenId string) (Token, error)
	ListByType(ctx context.Context, typ TokenType, limit uint64, offset uint64, order storage.SortOrder) ([]Token, error)
}

// Token -
type Token struct {
	bun.BaseModel `bun:"token"`

	ID          uint64          `bun:"id,pk,autoincrement" comment:"Unique internal identity"`
	FirstHeight uint64          `comment:"Block height when token was first time transferred or minted"`
	ContractId  uint64          `bun:",unique:token_unique_id" comment:"Token contract id"`
	TokenId     decimal.Decimal `bun:",unique:token_unique_id,type:numeric" comment:"Token id"`
	Type        TokenType       `bun:",type:token_type" comment:"Token type"`

	Contract Address `bun:"rel:belongs-to,join:contract_id=id" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (Token) TableName() string {
	return "token"
}

// GetHeight -
func (t Token) GetHeight() uint64 {
	return t.FirstHeight
}

// GetId -
func (t Token) GetId() uint64 {
	return t.ID
}

// String -
func (t Token) String() string {
	return fmt.Sprintf("%d_%s", t.ContractId, t.TokenId.String())
}
