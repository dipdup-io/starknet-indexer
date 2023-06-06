package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/shopspring/decimal"
)

// IInvoke -
type IInvoke interface {
	storage.Table[*Invoke]
	Copiable[Invoke]
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
	// nolint
	tableName struct{} `pg:"invoke,partition_by:RANGE(time),comment:Table with invokes"`

	ID                 uint64          `pg:"id,type:bigint,pk,notnull,comment:Unique internal identity"`
	Height             uint64          `pg:",use_zero,comment:Block height"`
	Time               time.Time       `pg:",pk,comment:Time of block"`
	Status             Status          `pg:",use_zero,comment:Status in blockchain (unknown - 1 | not received - 2  | received - 3 | pending - 4 | rejected - 5 | accepted on l2 - 6 | accepted on l1 - 7 )"`
	Hash               []byte          `pg:",comment:Transaction hash"`
	Version            uint64          `pg:",use_zero,comment:Version of invoke transaction"`
	Position           int             `pg:",use_zero,comment:Order in block"`
	ContractID         uint64          `pg:",comment:Contract address id"`
	EntrypointSelector []byte          `pg:",comment:Called selector"`
	Entrypoint         string          `pg:",comment:Entrypoint name"`
	MaxFee             decimal.Decimal `pg:",type:numeric,use_zero,comment:The maximum fee that the sender is willing to pay for the transaction"`
	Nonce              decimal.Decimal `pg:",type:numeric,use_zero,comment:The transaction nonce"`
	CallData           []string        `pg:",array,comment:Raw calldata"`
	ParsedCalldata     map[string]any  `pg:",comment:Calldata parsed according to contract ABI"`

	Contract  Address    `pg:"rel:has-one" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	Fee       *Fee       `pg:"rel:has-one"`
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
