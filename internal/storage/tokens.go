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

// IToken -
type IToken interface {
	storage.Table[*Token]

	ListByType(ctx context.Context, typ TokenType, limit uint64, offset uint64, order storage.SortOrder) ([]Token, error)
}

// Token -
type Token struct {
	// nolint
	tableName struct{} `pg:"token"`

	ID           uint64         `pg:",comment:Unique internal identity"`
	DeployHeight uint64         `pg:",use_zero,comment:Block height when token was deployed"`
	DeployTime   time.Time      `pg:",comment:Time of block when token was deployed"`
	ContractID   uint64         `pg:",comment:Token contract id"`
	OwnerID      uint64         `pg:",comment:Token owner id"`
	Type         TokenType      `pg:",type:SMALLINT,comment:Token type (1 - ERC20 | 2 - ERC721 | 3 - ERC1155)"`
	Metadata     map[string]any `pg:",type:jsonb,comment:Token metadata which was used as a constructor arguments"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `pg:"rel:has-one" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
}

// TableName -
func (Token) TableName() string {
	return "token"
}
