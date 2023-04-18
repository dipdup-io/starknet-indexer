package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IL1Handler -
type IL1Handler interface {
	storage.Table[*L1Handler]
	Copiable[L1Handler]
	Filterable[L1Handler, L1HandlerFilter]
}

// L1HandlerFilter -
type L1HandlerFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Contract       BytesFilter
	Selector       EqualityFilter
	Entrypoint     StringFilter
	ParsedCalldata map[string]string
}

// L1Handler -
type L1Handler struct {
	// nolint
	tableName struct{} `pg:"l1_handler,partition_by:RANGE(time),comment:Table with l1 handler transactions"`

	ID             uint64          `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height         uint64          `pg:",use_zero,comment:Block height"`
	Time           time.Time       `pg:",pk,comment:Time of block"`
	Status         Status          `pg:",use_zero,comment:Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	Hash           []byte          `pg:",comment:Transaction hash"`
	ContractID     uint64          `pg:",comment:Contract address id"`
	Position       int             `pg:",use_zero,comment:Order in block"`
	Selector       []byte          `pg:",comment:Called selector"`
	Entrypoint     string          `pg:",comment:Entrypoint name"`
	MaxFee         decimal.Decimal `pg:",type:numeric,use_zero,comment:The maximum fee that the sender is willing to pay for the transaction"`
	Nonce          decimal.Decimal `pg:",type:numeric,use_zero,comment:The transaction nonce"`
	CallData       []string        `pg:",array,comment:Raw calldata"`
	ParsedCalldata map[string]any  `pg:",comment:Calldata parsed according to contract ABI"`

	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
}

// TableName -
func (L1Handler) TableName() string {
	return "l1_handler"
}

// GetHeight -
func (l1 L1Handler) GetHeight() uint64 {
	return l1.Height
}
