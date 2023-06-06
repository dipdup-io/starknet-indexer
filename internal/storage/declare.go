package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
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
	// nolint
	tableName struct{} `pg:"declare,partition_by:RANGE(time),comment:Table with declare transactions"`

	ID         uint64          `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height     uint64          `pg:",use_zero,comment:Block height"`
	ClassID    uint64          `pg:",comment:Declared class id"`
	Version    uint64          `pg:",use_zero,comment:Declare transaction version"`
	Position   int             `pg:",use_zero,comment:Order in block"`
	SenderID   *uint64         `pg:",comment:Sender address id"`
	ContractID *uint64         `pg:",comment:Contract address id"`
	Time       time.Time       `pg:",pk,comment:Time of block"`
	Status     Status          `pg:",use_zero,comment:Status of block"`
	Hash       []byte          `pg:",comment:Transaction hash"`
	MaxFee     decimal.Decimal `pg:",type:numeric,use_zero,comment:The maximum fee that the sender is willing to pay for the transaction"`
	Nonce      decimal.Decimal `pg:",type:numeric,use_zero,comment:The transaction nonce"`

	Class     Class      `pg:"rel:has-one" hasura:"table:class,field:class_id,remote_field:id,type:oto,name:class"`
	Sender    Address    `pg:"rel:has-one" hasura:"table:address,field:sender_id,remote_field:id,type:oto,name:sender"`
	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
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
