package postgres

import (
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// ERC20 -
type ERC20 struct {
	*postgres.Table[*storage.ERC20]
}

// NewERC20 -
func NewERC20(db *database.PgGo) *ERC20 {
	return &ERC20{
		Table: postgres.NewTable[*storage.ERC20](db),
	}
}

// ERC721 -
type ERC721 struct {
	*postgres.Table[*storage.ERC721]
}

// NewERC721 -
func NewERC721(db *database.PgGo) *ERC721 {
	return &ERC721{
		Table: postgres.NewTable[*storage.ERC721](db),
	}
}

// ERC1155 -
type ERC1155 struct {
	*postgres.Table[*storage.ERC1155]
}

// NewERC1155 -
func NewERC1155(db *database.PgGo) *ERC1155 {
	return &ERC1155{
		Table: postgres.NewTable[*storage.ERC1155](db),
	}
}
