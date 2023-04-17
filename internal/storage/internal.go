package storage

import (
	"time"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

// IInternal -
type IInternal interface {
	storage.Table[*Internal]

	Copiable[Internal]
	Filterable[Internal, InternalFilter]
}

// InternalFilter -
type InternalFilter struct {
	ID             IntegerFilter
	Height         IntegerFilter
	Time           TimeFilter
	Status         EnumFilter
	Contract       BytesFilter
	Caller         BytesFilter
	Class          BytesFilter
	Selector       EqualityFilter
	Entrypoint     StringFilter
	EntrypointType EnumFilter
	CallType       EnumFilter
	ParsedCalldata map[string]string
}

// Internal -
type Internal struct {
	// nolint
	tableName struct{} `pg:"internal_tx,partition_by:RANGE(time)"`

	ID     uint64    `pg:"id,type:bigint,pk,notnull"`
	Height uint64    `pg:",use_zero"`
	Time   time.Time `pg:",pk"`

	DeclareID       *uint64
	DeployID        *uint64
	DeployAccountID *uint64
	InvokeID        *uint64
	L1HandlerID     *uint64
	InternalID      *uint64
	ClassID         uint64
	CallerID        uint64
	ContractID      uint64

	Status         Status         `pg:"type:SMALLINT,use_zero"`
	EntrypointType EntrypointType `pg:"type:SMALLINT"`
	CallType       CallType       `pg:"type:SMALLINT"`

	Hash           []byte
	Selector       []byte
	Entrypoint     string
	Result         []string `pg:",array"`
	Calldata       []string `pg:",array"`
	ParsedCalldata map[string]any
	ParsedResult   map[string]any

	Class     Class      `pg:"rel:has-one"`
	Caller    Address    `pg:"rel:has-one"`
	Contract  Address    `pg:"rel:has-one"`
	Internals []Internal `pg:"rel:has-many"`
	Messages  []Message  `pg:"rel:has-many"`
	Events    []Event    `pg:"rel:has-many"`
	Transfers []Transfer `pg:"rel:has-many"`
	ERC20     *ERC20     `pg:"-"`
	ERC721    *ERC721    `pg:"-"`
	ERC1155   *ERC1155   `pg:"-"`
}

// TableName -
func (Internal) TableName() string {
	return "internal_tx"
}

// GetHeight -
func (invoke Internal) GetHeight() uint64 {
	return invoke.Height
}
