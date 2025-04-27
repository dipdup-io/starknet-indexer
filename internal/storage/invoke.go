package storage

import (
	"github.com/dipdup-io/starknet-indexer/pkg/types"
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
	Countable[InvokeFilter]
	HashByHeight
}

// InvokeFilter -
type InvokeFilter struct {
	ID             IntegerFilter
	Hash           BytesFilter
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

	ID                 uint64          `bun:"id,type:bigint,pk,notnull" json:"id" comment:"Unique internal identity"`
	Height             uint64          `json:"height" comment:"Block height"`
	Time               time.Time       `bun:",pk" json:"time" comment:"Time of block"`
	Status             Status          `json:"status" comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	Hash               types.Hex       `json:"hash" comment:"Transaction hash"`
	Version            uint64          `json:"version" comment:"Version of invoke transaction"`
	Position           int             `json:"position" comment:"Order in block"`
	ContractID         uint64          `json:"contract_id" comment:"Contract address id"`
	EntrypointSelector types.Hex       `json:"entrypoint_selector" comment:"Called selector"`
	Entrypoint         string          `json:"entrypoint" comment:"Entrypoint name"`
	MaxFee             decimal.Decimal `bun:",type:numeric" json:"max_fee" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce              decimal.Decimal `bun:",type:numeric" json:"nonce" comment:"The transaction nonce"`
	CallData           []string        `bun:",array" json:"call_data" comment:"Raw calldata"`
	ParsedCalldata     map[string]any  `bun:",nullzero" json:"parsed_calldata" comment:"Calldata parsed according to contract ABI"`
	Error              *string         `bun:"error" json:"error" comment:"Reverted error"`

	Contract  Address    `bun:"rel:belongs-to" json:"contract" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many" json:"internals"`
	Messages  []Message  `bun:"rel:has-many" json:"messages"`
	Events    []Event    `bun:"rel:has-many" json:"events"`
	Transfers []Transfer `bun:"rel:has-many" json:"transfers"`
	Fee       *Fee       `bun:"rel:belongs-to" json:"fee"`
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
