package storage

import (
	"encoding/json"
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IDeclare interface {
	storage.Table[*Declare]
	Filterable[Declare, DeclareFilter]
	HashByHeight
}

// DeclareFilter -
type DeclareFilter struct {
	ID      IntegerFilter
	Height  IntegerFilter
	Time    TimeFilter
	Hash    BytesFilter
	Status  EnumFilter
	Version EnumFilter
}

// Declare -
type Declare struct {
	bun.BaseModel `bun:"declare" comment:"Table with declare transactions" partition:"RANGE(time)"`

	ID         uint64          `bun:"id,type:bigint,pk,notnull,nullzero" json:"id" comment:"Unique internal identity"`
	Height     uint64          `json:"height" comment:"Block height"`
	ClassID    uint64          `json:"class_id" comment:"Declared class id"`
	Version    uint64          `json:"version" comment:"Declare transaction version"`
	Position   int             `json:"position" comment:"Order in block"`
	SenderID   *uint64         `json:"sender_id" comment:"Sender address id"`
	ContractID *uint64         `json:"contract_id" comment:"Contract address id"`
	Time       time.Time       `bun:",pk" json:"time" comment:"Time of block"`
	Status     Status          `json:"status" comment:"Status of block"`
	Hash       []byte          `json:"hash" comment:"Transaction hash"`
	MaxFee     decimal.Decimal `bun:",type:numeric" json:"max_fee" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce      decimal.Decimal `bun:",type:numeric" json:"nonce" comment:"The transaction nonce"`
	Error      *string         `bun:"error" json:"error" comment:"Reverted error"`

	Class     Class      `bun:"rel:belongs-to" json:"class" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Sender    Address    `bun:"rel:belongs-to" json:"sender" hasura:"table:address,field:sender_id,remote_field:id,type:oto,name:sender"`
	Contract  Address    `bun:"rel:belongs-to" json:"contract" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many" json:"internals"`
	Messages  []Message  `bun:"rel:has-many" json:"messages"`
	Events    []Event    `bun:"rel:has-many" json:"events"`
	Transfers []Transfer `bun:"rel:has-many" json:"transfers"`
	Fee       *Fee       `bun:"rel:belongs-to" json:"fee"`
}

// TableName -
func (Declare) TableName() string {
	return "declare"
}

// GetHeight -
func (d Declare) GetHeight() uint64 {
	return d.Height
}

// GetId -
func (d Declare) GetId() uint64 {
	return d.ID
}

// Columns -
func (Declare) Columns() []string {
	return []string{
		"id", "height", "class_id", "version", "position",
		"sender_id", "contract_id", "time", "status", "hash",
		"max_fee", "nonce", "error",
	}
}

// Flat -
func (d Declare) Flat() []any {
	return []any{
		d.ID,
		d.Height,
		d.ClassID,
		d.Version,
		d.Position,
		d.SenderID,
		d.ContractID,
		d.Time,
		d.Status,
		d.Hash,
		d.MaxFee,
		d.Nonce,
		d.Error,
	}
}

func (d Declare) MarshalJSON() ([]byte, error) {
	type Alias Declare

	return json.Marshal(&struct {
		*Alias `json:"-"`
		Hash   string `json:"hash"`
	}{
		Alias: (*Alias)(&d),
		Hash:  BytesToFormattedHex(d.Hash),
	})
}
