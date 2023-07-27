package storage

import (
	"context"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// TokenType -
type TokenType int

// token types
const (
	TokenTypeERC20 TokenType = iota + 1
	TokenTypeERC721
	TokenTypeERC1155
)

// TokenFilter -
type TokenFilter struct {
	ID       IntegerFilter
	Contract BytesFilter
	Owner    BytesFilter
	Type     EnumFilter
}

// IToken -
type IToken interface {
	storage.Table[*Token]
	Filterable[Token, TokenFilter]

	ListByType(ctx context.Context, typ TokenType, limit uint64, offset uint64, order storage.SortOrder) ([]Token, error)
}

// Token -
type Token struct {
	// nolint
	tableName struct{} `pg:"token"`

	ID           uint64         `comment:"Unique internal identity"`
	DeployHeight uint64         `pg:",use_zero" comment:"Block height when token was deployed"`
	DeployTime   time.Time      `comment:"Time of block when token was deployed"`
	ContractID   uint64         `comment:"Token contract id"`
	OwnerID      uint64         `comment:"Token owner id"`
	Type         TokenType      `pg:",type:SMALLINT" comment:"Token type (1 - ERC20 | 2 - ERC721 | 3 - ERC1155)"`
	Metadata     map[string]any `pg:",type:jsonb" comment:"Token metadata which was used as a constructor arguments"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `pg:"rel:has-one" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
}

// TableName -
func (Token) TableName() string {
	return "token"
}

// GetHeight -
func (t Token) GetHeight() uint64 {
	return t.DeployHeight
}

// GetId -
func (t Token) GetId() uint64 {
	return t.ID
}
