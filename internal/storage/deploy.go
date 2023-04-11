package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IDeploy -
type IDeploy interface {
	storage.Table[*Deploy]
	Copiable[Deploy]
	Filterable[Deploy, DeployFilter]
}

// DeployFilter -
type DeployFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Class          BytesFilter
	ParsedCalldata map[string]string
}

// Deploy -
type Deploy struct {
	// nolint
	tableName struct{} `pg:"deploy,partition_by:RANGE(time)"`

	ID                  uint64 `pg:"id,type:bigint,pk,notnull"`
	Height              uint64 `pg:",use_zero"`
	ClassID             uint64
	ContractID          uint64
	Position            int       `pg:",use_zero"`
	Time                time.Time `pg:",pk"`
	Status              Status    `pg:",use_zero"`
	Hash                []byte
	ContractAddressSalt []byte
	ConstructorCalldata []string `pg:",array"`
	ParsedCalldata      map[string]any

	Class     Class      `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
	ERC20     *ERC20     `pg:"-"`
	ERC721    *ERC721    `pg:"-"`
	ERC1155   *ERC1155   `pg:"-"`
}

// TableName -
func (Deploy) TableName() string {
	return "deploy"
}

// GetHeight -
func (d Deploy) GetHeight() uint64 {
	return d.Height
}
