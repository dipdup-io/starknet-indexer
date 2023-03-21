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

	ID           uint64
	DeployHeight uint64 `pg:",use_zero"`
	DeployTime   time.Time
	ContractID   uint64
	Name         string
	Symbol       string
	Decimals     uint64

	Contract Address `pg:"rel:has-one"`
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

	ID           uint64
	DeployHeight uint64 `pg:",use_zero"`
	DeployTime   time.Time
	ContractID   uint64
	Name         string
	Symbol       string
	OwnerID      uint64

	Contract Address `pg:"rel:has-one"`
	Owner    Address `pg:"rel:has-one"`
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

	ID           uint64
	DeployHeight uint64 `pg:",use_zero"`
	DeployTime   time.Time
	ContractID   uint64
	TokenUri     string
	OwnerID      uint64

	Contract Address `pg:"rel:has-one"`
	Owner    Address `pg:"rel:has-one"`
}

// TableName -
func (ERC1155) TableName() string {
	return "erc1155"
}
