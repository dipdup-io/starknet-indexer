package api

type SqdBlockResponse struct {
	Header       BlockHeader     `json:"header"`
	Transactions []Transaction   `json:"transactions,omitempty"`
	Traces       []TraceResponse `json:"traces,omitempty"`
	Messages     []Message       `json:"messages,omitempty"`
	StateUpdates []StateUpdate   `json:"state_updates,omitempty"`
	StorageDiffs []StorageDiff   `json:"storage_diffs,omitempty"`
}

type BlockHeader struct {
	Number           uint64 `example:"321"                                                               json:"number"`
	Hash             string `example:"0x44529f2c44d9113e0ba4e53cb6e84f425ec186cda27545827b5a72d5540bfdc" json:"hash"`
	ParentHash       string `example:"0x44529f2c44d9113e0ba4e53cb6e84f425ec186cda27545827b5a72d5540bfdc" json:"parentHash"`
	Status           string `example:"ACCEPTED_ON_L1"                                                    json:"status"`
	NewRoot          string `example:"0x44529f2c44d9113e0ba4e53cb6e84f425ec186cda27545827b5a72d5540bfdc" json:"newRoot"`
	Timestamp        int64  `example:"1641950335"                                                        json:"timestamp"`
	SequencerAddress string `example:"0x44529f2c44d9113e0ba4e53cb6e84f425ec186cda27545827b5a72d5540bfdc" json:"sequencerAddress"`
}

type Transaction struct {
	TransactionIndex    uint      `example:"0"                                                                 json:"transactionIndex"`
	TransactionHash     string    `example:"0x794fae89c8c4b8f5f77a4996948d2547740f90e54bb4a5cc6119a7c70eca42c" json:"transactionHash"`
	ContractAddress     *string   `example:"0x1cee8364383aea317eefc181dbd8732f1504fd4511aed58f32c369dd546da0d" json:"contractAddress"`
	EntryPointSelector  *string   `example:"0x317eb442b72a9fae758d4fb26830ed0d9f31c8e7da4dbff4e8c59ea6a158e7f" json:"entryPointSelector"`
	Calldata            *[]string `json:"calldata"`
	MaxFee              *string   `example:"0x0"                                                               json:"maxFee"`
	Type                string    `example:"INVOKE"                                                            json:"type"`
	SenderAddress       *string   `json:"senderAddress"`
	Version             string    `example:"0x0"                                                               json:"version"`
	Signature           *[]string `json:"signature"`
	Nonce               *uint64   `json:"nonce"`
	ClassHash           *string   `json:"classHash"`
	CompiledClassHash   *string   `json:"compiledClassHash"`
	ContractAddressSalt *string   `json:"contractAddressSalt"`
	ConstructorCalldata *[]string `json:"constructorCalldata"`
}

type TraceResponse struct {
	TransactionIndex   uint     `json:"transaction_index"`
	TraceAddress       []int    `json:"trace_address"`
	TraceType          string   `json:"traceType"`
	InvocationType     string   `json:"invocationType"`
	CallerAddress      string   `json:"callerAddress"`
	ContractAddress    string   `json:"contractAddress"`
	CallType           *string  `json:"callType"`
	ClassHash          *string  `json:"classHash"`
	EntryPointSelector *string  `json:"entryPointSelector"`
	EntryPointType     *string  `json:"entryPointType"`
	Calldata           []string `json:"calldata"`
	Result             []string `json:"result"`
}

type Message struct {
	TransactionIndex uint     `json:"transaction_index"`
	TraceAddress     []int    `json:"trace_address"`
	Order            uint     `json:"order"`
	FromAddress      *string  `json:"fromAddress"`
	ToAddress        string   `json:"toAddress"`
	Payload          []string `json:"payload"`
}

type StateUpdate struct {
	NewRoot           string             `json:"newRoot"`
	OldRoot           string             `json:"oldRoot"`
	DeprecatedClasses []any              `json:"deprecatedDeclaredClasses"`
	DeclaredClasses   []any              `json:"declaredClasses"`
	DeployedContracts []DeployedContract `json:"deployedContracts"`
	ReplacedClasses   []any              `json:"replacedClasses"`
	Nonces            []any              `json:"nonces"`
}

type DeployedContract struct {
	Address   string `json:"address"`
	ClassHash string `json:"class_hash"`
}

type StorageDiff struct {
	Address string `json:"address"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}
