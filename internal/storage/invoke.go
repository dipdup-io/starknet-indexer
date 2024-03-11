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
type IInvoke interface {
	storage.Table[*Invoke]
	Filterable[Invoke, InvokeFilter]
}

// InvokeFilter -
type InvokeFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Version        EnumFilter
	Contract       BytesFilter
	Selector       EqualityFilter
	Entrypoint     StringFilter
	ParsedCalldata map[string]string
}

// Invoke -
type Invoke struct {
	bun.BaseModel `bun:"invoke" comment:"Table with invokes" partition:"RANGE(time)"`

	ID                 uint64          `bun:"id,type:bigint,pk,notnull" comment:"Unique internal identity"`
	Height             uint64          `comment:"Block height"`
	Time               time.Time       `bun:",pk" comment:"Time of block"`
	Status             Status          `comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	Hash               []byte          `comment:"Transaction hash"`
	Version            uint64          `comment:"Version of invoke transaction"`
	Position           int             `comment:"Order in block"`
	ContractID         uint64          `comment:"Contract address id"`
	EntrypointSelector []byte          `comment:"Called selector"`
	Entrypoint         string          `comment:"Entrypoint name"`
	MaxFee             decimal.Decimal `bun:",type:numeric" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce              decimal.Decimal `bun:",type:numeric" comment:"The transaction nonce"`
	CallData           []string        `bun:",array" comment:"Raw calldata"`
	ParsedCalldata     map[string]any  `bun:",nullzero" comment:"Calldata parsed according to contract ABI"`
	Error              *string         `bun:"error" comment:"Reverted error"`

	Contract  Address    `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many"`
	Messages  []Message  `bun:"rel:has-many"`
	Events    []Event    `bun:"rel:has-many"`
	Transfers []Transfer `bun:"rel:has-many"`
	Fee       *Fee       `bun:"rel:belongs-to"`
}

// TableName -
func (Invoke) TableName() string {
	return "invoke"
}

// GetHeight -
func (invoke Invoke) GetHeight() uint64 {
	return invoke.Height
}

// GetId -
func (invoke Invoke) GetId() uint64 {
	return invoke.ID
}

// Columns -
func (Invoke) Columns() []string {
	return []string{
		"id", "height", "time", "status", "hash", "version",
		"position", "contract_id", "entrypoint_selector",
		"entrypoint", "max_fee", "nonce", "call_data", "parsed_calldata",
		"error",
	}
}

// Flat -
func (i Invoke) Flat() []any {
	data := []any{
		i.ID,
		i.Height,
		i.Time,
		i.Status,
		i.Hash,
		i.Version,
		i.Position,
		i.ContractID,
		i.EntrypointSelector,
		i.Entrypoint,
		i.MaxFee,
		i.Nonce,
		pq.StringArray(i.CallData),
		nil,
		i.Error,
	}
	if i.ParsedCalldata != nil {
		parsed, err := json.MarshalWithOption(i.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
		if err == nil {
			data[13] = string(parsed)
		}
	}

	return data
}
