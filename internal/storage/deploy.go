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
	tableName struct{} `pg:"deploy,partition_by:RANGE(time)" comment:"Table with deploy transactions"`

	ID                  uint64         `pg:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height              uint64         `pg:",use_zero" comment:"Block height"`
	ClassID             uint64         `comment:"Class id"`
	ContractID          uint64         `comment:"Contract address id"`
	Position            int            `pg:",use_zero" comment:"Order in block"`
	Time                time.Time      `pg:",pk" comment:"Time of block"`
	Status              Status         `pg:",use_zero"`
	Hash                []byte         `comment:"Transaction hash"`
	ContractAddressSalt []byte         `comment:"A random salt that determines the account address"`
	ConstructorCalldata []string       `pg:",array" comment:"Raw constructor calldata"`
	ParsedCalldata      map[string]any `comment:"Calldata parsed according to contract ABI"`

	Class     Class      `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
	Token     *Token     `pg:"-"`
}

// TableName -
func (Deploy) TableName() string {
	return "deploy"
}

// GetHeight -
func (d Deploy) GetHeight() uint64 {
	return d.Height
}

// GetId -
func (d Deploy) GetId() uint64 {
	return d.ID
}
