package api

type Request struct {
	Type             string                 `json:"type"`
	FromBlock        uint64                 `json:"fromBlock"`
	ToBlock          uint64                 `json:"toBlock,omitempty"`
	IncludeAllBlocks bool                   `json:"includeAllBlocks"`
	Fields           Fields                 `json:"fields,omitempty"`
	StateUpdates     []map[string]any       `json:"stateUpdates,omitempty"`
	StorageDiffs     []map[string]any       `json:"storageDiffs,omitempty"`
	Traces           []Trace                `json:"traces,omitempty"`
	Messages         []map[string]any       `json:"messages,omitempty"`
	Transactions     []TransactionWithTrace `json:"transactions,omitempty"`
}

type Fields struct {
	Block       BlockField       `json:"block"`
	StateUpdate StateUpdateField `json:"stateUpdate"`
	StorageDiff StorageDiffField `json:"storageDiff"`
	Trace       TraceField       `json:"trace"`
	Transaction TransactionField `json:"transaction"`
	Event       EventField       `json:"event"`
	Message     MessageField     `json:"message"`
}

type BlockField struct {
	Timestamp        bool `json:"timestamp"`
	ParentHash       bool `json:"parentHash,omitempty"`
	Status           bool `json:"status,omitempty"`
	NewRoot          bool `json:"newRoot,omitempty"`
	SequencerAddress bool `json:"sequencerAddress,omitempty"`
}

type StateUpdateField struct {
	NewRoot                   bool `json:"newRoot"`
	OldRoot                   bool `json:"oldRoot"`
	DeprecatedDeclaredClasses bool `json:"deprecatedDeclaredClasses"`
	DeclaredClasses           bool `json:"declaredClasses"`
	DeployedContracts         bool `json:"deployedContracts"`
	ReplacedClasses           bool `json:"replacedClasses"`
	Nonces                    bool `json:"nonces"`
}

type StorageDiffField struct {
	Value bool `json:"value"`
}

type TraceField struct {
	TraceType          bool `json:"traceType"`
	InvocationType     bool `json:"invocationType"`
	CallerAddress      bool `json:"callerAddress"`
	ContractAddress    bool `json:"contractAddress"`
	CallType           bool `json:"callType"`
	ClassHash          bool `json:"classHash"`
	EntryPointSelector bool `json:"entryPointSelector"`
	EntryPointType     bool `json:"entryPointType"`
	Calldata           bool `json:"calldata"`
	Result             bool `json:"result"`
}

type TransactionField struct {
	TransactionHash     bool `json:"transactionHash"`
	ContractAddress     bool `json:"contractAddress"`
	EntryPointSelector  bool `json:"entryPointSelector"`
	Calldata            bool `json:"calldata"`
	MaxFee              bool `json:"maxFee"`
	Type                bool `json:"type"`
	SenderAddress       bool `json:"senderAddress"`
	Version             bool `json:"version"`
	Signature           bool `json:"signature"`
	Nonce               bool `json:"nonce"`
	ClassHash           bool `json:"classHash"`
	CompiledClassHash   bool `json:"compiledClassHash"`
	ContractAddressSalt bool `json:"contractAddressSalt"`
	ConstructorCalldata bool `json:"constructorCalldata"`
}

type EventField struct {
	Keys         bool `json:"keys"`
	Data         bool `json:"data"`
	TraceAddress bool `json:"traceAddress"`
}

type MessageField struct {
	FromAddress bool `json:"fromAddress"`
	ToAddress   bool `json:"toAddress"`
	Payload     bool `json:"payload"`
}

type Trace struct {
	Events bool `json:"events"`
}

type TransactionWithTrace struct {
	Traces bool `json:"traces"`
	Events bool `json:"events"`
}

func NewRequest(fromLevel uint64, toLevel uint64) *Request {
	return &Request{
		Type:             "starknet",
		FromBlock:        fromLevel,
		ToBlock:          toLevel,
		IncludeAllBlocks: true,
		Fields: Fields{
			Block: BlockField{
				ParentHash:       true,
				Status:           true,
				NewRoot:          true,
				Timestamp:        true,
				SequencerAddress: true,
			},
			StateUpdate: StateUpdateField{
				NewRoot:                   true,
				OldRoot:                   true,
				DeprecatedDeclaredClasses: true,
				DeclaredClasses:           true,
				DeployedContracts:         true,
				ReplacedClasses:           true,
				Nonces:                    true,
			},
			StorageDiff: StorageDiffField{
				Value: true,
			},
			Trace: TraceField{
				TraceType:          true,
				InvocationType:     true,
				CallerAddress:      true,
				ContractAddress:    true,
				CallType:           true,
				ClassHash:          true,
				EntryPointSelector: true,
				EntryPointType:     true,
				Calldata:           true,
				Result:             true,
			},
			Transaction: TransactionField{
				TransactionHash:     true,
				ContractAddress:     true,
				EntryPointSelector:  true,
				Calldata:            true,
				MaxFee:              true,
				Type:                true,
				SenderAddress:       true,
				Version:             true,
				Signature:           true,
				Nonce:               true,
				ClassHash:           true,
				CompiledClassHash:   true,
				ContractAddressSalt: true,
				ConstructorCalldata: true,
			},
			Event: EventField{
				Keys:         true,
				Data:         true,
				TraceAddress: true,
			},
			Message: MessageField{
				FromAddress: true,
				ToAddress:   true,
				Payload:     true,
			},
		},
		StateUpdates: []map[string]any{
			{},
		},
		StorageDiffs: []map[string]any{
			{},
		},
		Traces: []Trace{
			{Events: true},
		},
		Messages: []map[string]any{
			{},
		},
		Transactions: []TransactionWithTrace{
			{Traces: true, Events: true},
		},
	}
}

func NewSimpleRequest(level uint64) *Request {
	return &Request{
		Type:             "starknet",
		FromBlock:        level,
		IncludeAllBlocks: true,
	}
}
