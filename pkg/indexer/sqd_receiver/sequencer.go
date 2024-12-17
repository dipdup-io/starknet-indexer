package sqd_receiver

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/sqd_receiver/api"
)

func (r *Receiver) sequencer(ctx context.Context) {
	orderedBlocks := map[uint64]*api.SqdBlockResponse{}
	l := r.Level()
	currentBlock := l + 1

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
				r.MustOutput(OutputName).Push(b)
				r.Log.Info().
					Uint64("ID", b.Header.Number).
					Msg("sended block")

				r.setLevel(currentBlock)
				r.Log.Debug().
					Uint64("height", currentBlock).
					Msg("put in order block")

				delete(orderedBlocks, currentBlock)
				currentBlock += 1

				b, ok = orderedBlocks[currentBlock]
			}
		}
	}
}
