package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// IDeclare -
type IDeclare interface {
	storage.Table[*Declare]
	Filterable[Declare, DeclareFilter]
}

// DeclareFilter -
type DeclareFilter struct {
	ID      IntegerFilter
	Height  IntegerFilter
	Time    TimeFilter
	Status  EnumFilter
	Version EnumFilter
}

// Declare -
type Declare struct {
	bun.BaseModel `bun:"declare" comment:"Table with declare transactions" partition:"RANGE(time)"`

	ID         uint64          `bun:"id,type:bigint,pk,notnull,nullzero" comment:"Unique internal identity"`
	Height     uint64          `comment:"Block height"`
	ClassID    uint64          `comment:"Declared class id"`
	Version    uint64          `comment:"Declare transaction version"`
	Position   int             `comment:"Order in block"`
	SenderID   *uint64         `comment:"Sender address id"`
	ContractID *uint64         `comment:"Contract address id"`
	Time       time.Time       `bun:",pk" comment:"Time of block"`
	Status     Status          `comment:"Status of block"`
	Hash       []byte          `comment:"Transaction hash"`
	MaxFee     decimal.Decimal `bun:",type:numeric" comment:"The maximum fee that the sender is willing to pay for the transaction"`
	Nonce      decimal.Decimal `bun:",type:numeric" comment:"The transaction nonce"`

	Class     Class      `bun:"rel:belongs-to" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Sender    Address    `bun:"rel:belongs-to" hasura:"table:address,field:sender_id,remote_field:id,type:oto,name:sender"`
	Contract  Address    `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `bun:"rel:has-many"`
	Messages  []Message  `bun:"rel:has-many"`
	Events    []Event    `bun:"rel:has-many"`
	Transfers []Transfer `bun:"rel:has-many"`
	Fee       *Fee       `bun:"rel:belongs-to"`
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
		"max_fee", "nonce",
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
	}
}
