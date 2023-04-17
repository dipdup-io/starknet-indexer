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
	tableName struct{} `pg:"deploy_account,partition_by:RANGE(time)"`

	ID                  uint64 `pg:"id,type:bigint,pk,notnull"`
	Height              uint64 `pg:",use_zero"`
	ClassID             uint64
	ContractID          uint64
	Position            int       `pg:",use_zero"`
	Time                time.Time `pg:",pk"`
	Status              Status    `pg:",use_zero"`
	Hash                []byte
	ContractAddressSalt []byte
	MaxFee              decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce               decimal.Decimal `pg:",type:numeric,use_zero"`
	ConstructorCalldata []string        `pg:",array"`
	ParsedCalldata      map[string]any

	Class     Class      `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
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
