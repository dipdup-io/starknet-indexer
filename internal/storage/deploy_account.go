package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IDeployAccount -
type IDeployAccount interface {
	storage.Table[*DeployAccount]
	Filterable[DeployAccount, DeployAccountFilter]
}

// DeployAccountFilter -
type DeployAccountFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Class          BytesFilter
	ParsedCalldata map[string]string
}

// DeployAccount -
type DeployAccount struct {
	// nolint
	tableName struct{} `pg:"deploy_account,partition_by:RANGE(time),comment:table with deploy account transactions"`

	ID                  uint64          `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height              uint64          `pg:",use_zero,comment:Block height"`
	ClassID             uint64          `pg:",comment:Class id"`
	ContractID          uint64          `pg:",comment:Contract address id"`
	Position            int             `pg:",use_zero,comment:Order in block"`
	Time                time.Time       `pg:",pk,comment:Time of block"`
	Status              Status          `pg:",use_zero"`
	Hash                []byte          `pg:",comment:Transaction hash"`
	ContractAddressSalt []byte          `pg:",comment:A random salt that determines the account address"`
	MaxFee              decimal.Decimal `pg:",type:numeric,use_zero,comment:The maximum fee that the sender is willing to pay for the transaction"`
	Nonce               decimal.Decimal `pg:",type:numeric,use_zero,comment:The transaction nonce"`
	ConstructorCalldata []string        `pg:",array,comment:Raw constructor calldata"`
	ParsedCalldata      map[string]any  `pg:",comment:Calldata parsed according to contract ABI"`

	Class     Class      `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
}

// TableName -
func (DeployAccount) TableName() string {
	return "deploy_account"
}

// GetHeight -
func (d DeployAccount) GetHeight() uint64 {
	return d.Height
}

// GetId -
func (d DeployAccount) GetId() uint64 {
	return d.ID
}
