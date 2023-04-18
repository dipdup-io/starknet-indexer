package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IERC20 -
type IERC20 interface {
	storage.Table[*ERC20]
}

// ERC20 -
type ERC20 struct {
	// nolint
	tableName struct{} `pg:"erc20"`

	ID           uint64    `pg:",comment:Unique internal identity"`
	DeployHeight uint64    `pg:",use_zero,comment:Block height when token was deployed"`
	DeployTime   time.Time `pg:",comment:Time of block when token was deployed"`
	ContractID   uint64    `pg:",comment:Token contract id"`
	Name         string    `pg:",comment:Token name"`
	Symbol       string    `pg:",comment:Token symbol"`
	Decimals     uint64    `pg:",comment:Token decimals"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
}

// TableName -
func (ERC20) TableName() string {
	return "erc20"
}

// IERC721 -
type IERC721 interface {
	storage.Table[*ERC721]
}

// ERC721 -
type ERC721 struct {
	// nolint
	tableName struct{} `pg:"erc721"`

	ID           uint64    `pg:",comment:Unique internal identity"`
	DeployHeight uint64    `pg:",use_zero,comment:Block height when token was deployed"`
	DeployTime   time.Time `pg:",comment:Time of block when token was deployed"`
	ContractID   uint64    `pg:",comment:Token contract id"`
	Name         string    `pg:",comment:Token name"`
	Symbol       string    `pg:",comment:Token symbol"`
	OwnerID      uint64    `pg:",comment:Token owner id"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `pg:"rel:has-one" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
}

// TableName -
func (ERC721) TableName() string {
	return "erc721"
}

// IERC1155 -
type IERC1155 interface {
	storage.Table[*ERC1155]
}

// ERC1155 -
type ERC1155 struct {
	// nolint
	tableName struct{} `pg:"erc1155"`

	ID           uint64    `pg:",comment:Unique internal identity"`
	DeployHeight uint64    `pg:",use_zero,comment:Block height when token was deployed"`
	DeployTime   time.Time `pg:",comment:Time of block when token was deployed"`
	ContractID   uint64    `pg:",comment:Token contract id"`
	TokenUri     string    `pg:",comment:Token uri"`
	OwnerID      uint64    `pg:",comment:Token owner id"`

	Contract Address `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Owner    Address `pg:"rel:has-one" hasura:"table:address,field:owner_id,remote_field:id,type:oto,name:owner"`
}

// TableName -
func (ERC1155) TableName() string {
	return "erc1155"
}
