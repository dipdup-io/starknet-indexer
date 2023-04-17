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
	tableName struct{} `pg:"invoke,partition_by:RANGE(time)"`

	ID                 uint64    `pg:"id,type:bigint,pk,notnull"`
	Height             uint64    `pg:",use_zero"`
	Time               time.Time `pg:",pk"`
	Status             Status    `pg:",use_zero"`
	Hash               []byte
	Version            uint64 `pg:",use_zero"`
	Position           int    `pg:",use_zero"`
	ContractID         uint64
	EntrypointSelector []byte
	Entrypoint         string
	MaxFee             decimal.Decimal `pg:",type:numeric,use_zero"`
	Nonce              decimal.Decimal `pg:",type:numeric,use_zero"`
	CallData           []string        `pg:",array"`
	ParsedCalldata     map[string]any

	Contract  Address    `pg:"rel:has-one"`
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
