package receiver

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
)

func (r *Receiver) sequencer(ctx context.Context) {
	orderedBlocks := map[uint64]*api.SqdBlockResponse{}
	currentBlock := r.Level()
	if currentBlock != 0 {
		currentBlock += 1
	}

	for {
		select {
		case <-ctx.Done():
			return
		case block, ok := <-r.blocks:
			if !ok {
				r.Log.Warn().Msg("can't read message from input, it was drained and closed")
				r.MustOutput(StopOutput).Push(struct{}{})
				return
			}
			orderedBlocks[block.Header.Number] = block

			b, ok := orderedBlocks[currentBlock]
			for ok {
				r.MustOutput(BlocksOutput).Push(b)
				r.setLevel(currentBlock)
				delete(orderedBlocks, currentBlock)
				currentBlock += 1

				b, ok = orderedBlocks[currentBlock]
			}
		}
	}
}
