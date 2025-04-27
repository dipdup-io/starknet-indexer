package block

import (
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/types"
	"github.com/goccy/go-json"
	"time"
)

type Block struct {
	Height             uint64    `json:"height"`
	Time               time.Time `json:"time"`
	Version            *string   `json:"version"`
	TxCount            int       `json:"tx_count"`
	InvokeCount        int       `json:"invoke_count"`
	DeclareCount       int       `json:"declare_count"`
	DeployCount        int       `json:"deploy_count"`
	DeployAccountCount int       `json:"deploy_account_count"`
	L1HandlerCount     int       `json:"l1_handler_count"`
	StorageDiffCount   int       `json:"storage_diff_count"`
	Status             string    `json:"status"`
	Hash               types.Hex `json:"hash"`
	ParentHash         types.Hex `json:"parent_hash"`
	NewRoot            types.Hex `json:"new_root"`
	SequencerAddress   types.Hex `json:"sequencer_address"`
}

func marshalBlock(block models.Block) ([]byte, error) {
	return json.Marshal(Block{
		Height:             block.Height,
		Time:               block.Time,
		Version:            block.Version,
		TxCount:            block.TxCount,
		InvokeCount:        block.InvokeCount,
		DeclareCount:       block.DeclareCount,
		DeployCount:        block.DeployCount,
		DeployAccountCount: block.DeployAccountCount,
		L1HandlerCount:     block.L1HandlerCount,
		StorageDiffCount:   block.StorageDiffCount,
		Status:             block.Status.String(),
		Hash:               block.Hash,
		ParentHash:         block.ParentHash,
		NewRoot:            block.NewRoot,
		SequencerAddress:   block.SequencerAddress,
	})
}
