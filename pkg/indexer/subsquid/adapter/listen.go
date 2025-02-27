package adapter

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
)

func (a *Adapter) listen(ctx context.Context) {
	a.Log.Info().Msg("module started")

	blocksInput := a.MustInput(BlocksInput)
	headInput := a.MustInput(HeadInput)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-blocksInput.Listen():
			if !ok {
				a.Log.Warn().Msg("can't read message from input, it was drained and closed")
				a.MustOutput(StopOutput).Push(struct{}{})
				return
			}

			block, ok := msg.(*api.SqdBlockResponse)

			if !ok {
				a.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			results, err := a.convert(ctx, block)
			if err != nil {
				a.Log.Err(err).
					Uint64("height", block.Header.Number).
					Msg("convert error")
				a.MustOutput(StopOutput).Push(struct{}{})
				continue
			}
			a.results <- results
			if results.Block.Height == a.head {
				a.MustOutput(HeadAchieved).Push(struct{}{})
				return
			}

		case msg, ok := <-headInput.Listen():
			if !ok {
				a.Log.Warn().Msg("can't read message from input, it was drained and closed")
				a.MustOutput(StopOutput).Push(struct{}{})
				return
			}
			head, ok := msg.(uint64)
			if !ok {
				a.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}
			a.head = head
		}
	}
}
