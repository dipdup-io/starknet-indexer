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
type IL1Handler interface {
	storage.Table[*L1Handler]
	Filterable[L1Handler, L1HandlerFilter]
	HashByHeight
}

// L1HandlerFilter -
type L1HandlerFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Contract       BytesFilter
	Hash           BytesFilter
	Selector       EqualityFilter
	Entrypoint     StringFilter
	ParsedCalldata map[string]string
}

// L1Handler -
type L1Handler struct {
	bun.BaseModel `bun:"l1_handler" comment:"Table with l1 handler transactions" partition:"RANGE(time)"`

	ID             uint64          `bun:"id,type:bigint,pk,notnull" json:"id" comment:"Unique internal identity"`
	Height         uint64          `json:"height" comment:"Block height"`
	Time           time.Time       `bun:",pk" json:"time" comment:"Time of block"`
	Status         Status          `json:"status" comment:"Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 | reverted - 8)"`
	Hash           types.Hex       `json:"hash" comment:"Transaction hash"`
	ContractID     uint64          `json:"contract_id" comment:"Contract address id"`
	Position       int             `json:"position" comment:"Order in block"`
	Selector       types.Hex       `json:"selector" comment:"Called selector"`
	Entrypoint     string          `json:"entrypoint" comment:"Entrypoint name"`
	MaxFee         decimal.Decimal `bun:",type:numeric" json:"max_fee" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce          decimal.Decimal `bun:",type:numeric" json:"nonce" comment:"The transaction nonce"`
	CallData       []string        `bun:",array" json:"call_data" comment:"Raw calldata"`
	ParsedCalldata map[string]any  `bun:",nullzero" json:"parsed_calldata" comment:"Calldata parsed according to contract ABI"`
	Error          *string         `bun:"error" json:"error" comment:"Reverted error"`

	Contract  Address    `bun:"rel:belongs-to" json:"contract" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many" json:"internals"`
	Messages  []Message  `bun:"rel:has-many" json:"messages"`
	Transfers []Transfer `bun:"rel:has-many" json:"transfers"`
	Events    []Event    `bun:"rel:has-many" json:"events"`
	Fee       *Fee       `bun:"rel:belongs-to" json:"fee"`
}

// TableName -
func (L1Handler) TableName() string {
	return "l1_handler"
}

// GetHeight -
func (l1 L1Handler) GetHeight() uint64 {
	return l1.Height
}

// GetId -
func (l1 L1Handler) GetId() uint64 {
	return l1.ID
}

// Columns -
func (L1Handler) Columns() []string {
	return []string{
		"id", "height", "time", "status", "hash", "contract_id",
		"position", "selector", "entrypoint", "max_fee",
		"nonce", "call_data", "parsed_calldata", "error",
	}
}

// Flat -
func (l1 L1Handler) Flat() []any {
	data := []any{
		l1.ID,
		l1.Height,
		l1.Time,
		l1.Status,
		l1.Hash,
		l1.ContractID,
		l1.Position,
		l1.Selector,
		l1.Entrypoint,
		l1.MaxFee,
		l1.Nonce,
		pq.StringArray(l1.CallData),
		nil,
		l1.Error,
	}

	if l1.ParsedCalldata != nil {
		parsed, err := json.MarshalWithOption(l1.ParsedCalldata, json.UnorderedMap(), json.DisableNormalizeUTF8())
		if err == nil {
			data[12] = string(parsed)
		}
	}
	return data
}

//func (h L1Handler) MarshalJSON() ([]byte, error) {
//	type Alias L1Handler
//
//	return json.Marshal(&struct {
//		Alias
//		Hash     string `json:"hash"`
//		Selector string `json:"selector"`
//	}{
//		Alias:    Alias(h),
//		Hash:     BytesToFormattedHex(h.Hash),
//		Selector: BytesToFormattedHex(h.Selector),
//	})
//}
