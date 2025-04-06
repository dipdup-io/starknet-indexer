package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IDeployAccount interface {
	storage.Table[*DeployAccount]
	Filterable[DeployAccount, DeployAccountFilter]
	HashByHeight
}

// DeployAccountFilter -
type DeployAccountFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Class          BytesFilter
	Hash           BytesFilter
	ParsedCalldata map[string]string
}

// DeployAccount -
type DeployAccount struct {
	bun.BaseModel `bun:"deploy_account" comment:"table with deploy account transactions" partition:"RANGE(time)"`

	ID                  uint64          `bun:"id,type:bigint,pk,notnull" json:"id" comment:"Unique internal identity"`
	Height              uint64          `json:"height" comment:"Block height"`
	ClassID             uint64          `json:"class_id" comment:"Class id"`
	ContractID          uint64          `json:"contract_id" comment:"Contract address id"`
	Position            int             `json:"position" comment:"Order in block"`
	Time                time.Time       `bun:",pk" json:"time" comment:"Time of block"`
	Status              Status          `json:"status"`
	Hash                []byte          `json:"hash" comment:"Transaction hash"`
	ContractAddressSalt []byte          `json:"contract_address_salt" comment:"A random salt that determines the account address"`
	MaxFee              decimal.Decimal `bun:",type:numeric" json:"max_fee" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce               decimal.Decimal `bun:",type:numeric" json:"nonce" comment:"The transaction nonce"`
	ConstructorCalldata []string        `bun:",array" json:"constructor_calldata" comment:"Raw constructor calldata"`
	ParsedCalldata      map[string]any  `bun:",nullzero" json:"parsed_calldata" comment:"Calldata parsed according to contract ABI"`
	Error               *string         `bun:"error" json:"error" comment:"Reverted error"`

	Class     Class      `bun:"rel:belongs-to" json:"class" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Contract  Address    `bun:"rel:belongs-to" json:"contract" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many" json:"internals"`
	Messages  []Message  `bun:"rel:has-many" json:"messages"`
	Events    []Event    `bun:"rel:has-many" json:"events"`
	Transfers []Transfer `bun:"rel:has-many" json:"transfers"`
	Fee       *Fee       `bun:"rel:belongs-to" json:"fee"`
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
		nil,
		d.Error,
	}
	if d.ParsedCalldata != nil {
		parsed, err := json.MarshalWithOption(d.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
		if err == nil {
			data[12] = string(parsed)
		}
	}
	return data
}

func (d DeployAccount) MarshalJSON() ([]byte, error) {
	type Alias DeployAccount

	return json.Marshal(&struct {
		Alias
		Hash                string `json:"hash"`
		ContractAddressSalt string `json:"contract_address_salt"`
	}{
		Alias:               Alias(d),
		Hash:                BytesToFormattedHex(d.Hash),
		ContractAddressSalt: BytesToFormattedHex(d.ContractAddressSalt),
	})
}
