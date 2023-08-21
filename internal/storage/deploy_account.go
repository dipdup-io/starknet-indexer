package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
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
	bun.BaseModel `bun:"deploy_account" comment:"table with deploy account transactions" partition:"RANGE(time)"`

	ID                  uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height              uint64          `comment:"Block height"`
	ClassID             uint64          `comment:"Class id"`
	ContractID          uint64          `comment:"Contract address id"`
	Position            int             `comment:"Order in block"`
	Time                time.Time       `bun:",pk" comment:"Time of block"`
	Status              Status          ``
	Hash                []byte          `comment:"Transaction hash"`
	ContractAddressSalt []byte          `comment:"A random salt that determines the account address"`
	MaxFee              decimal.Decimal `bun:",type:numeric" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce               decimal.Decimal `bun:",type:numeric" comment:"The transaction nonce"`
	ConstructorCalldata []string        `bun:",array" comment:"Raw constructor calldata"`
	ParsedCalldata      map[string]any  `comment:"Calldata parsed according to contract ABI"`
	Error               *string         `bun:"error" comment:"Reverted error"`

	Class     Class      `bun:"rel:belongs-to" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Contract  Address    `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many"`
	Messages  []Message  `bun:"rel:has-many"`
	Events    []Event    `bun:"rel:has-many"`
	Transfers []Transfer `bun:"rel:has-many"`
	Fee       *Fee       `bun:"rel:belongs-to"`
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

// Columns -
func (DeployAccount) Columns() []string {
	return []string{
		"id", "height", "class_id", "contract_id", "position",
		"time", "status", "hash", "contract_address_salt",
		"max_fee", "nonce", "constructor_calldata", "parsed_calldata",
		"error",
	}
}

// Flat -
func (d DeployAccount) Flat() []any {
	data := []any{
		d.ID,
		d.Height,
		d.ClassID,
		d.ContractID,
		d.Position,
		d.Time,
		d.Status,
		d.Hash,
		d.ContractAddressSalt,
		d.MaxFee,
		d.Nonce,
		pq.StringArray(d.ConstructorCalldata),
		d.Error,
	}
	parsed, err := json.MarshalWithOption(d.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
	if err != nil {
		data = append(data, nil)
	} else {
		data = append(data, string(parsed))
	}
	return data
}
