package sqd_receiver

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (r *Receiver) worker(ctx context.Context, blockRange BlocksToWorker) {
	from := blockRange.From
	var allBlocksDownloaded bool

	for !allBlocksDownloaded {
		select {
		case <-ctx.Done():
			return
		default:
			blocks, err := r.api.GetBlocks(ctx, from, blockRange.To, blockRange.WorkerURL)
			if err != nil {
				log.Err(err).
					Uint64("fromLevel", from).
					Uint64("toLevel", blockRange.To).
					Str("worker url", blockRange.WorkerURL).
					Msg("loading blocks error")
				return
			}

			lastBlock := blocks[len(blocks)-1]

			for _, block := range blocks {
				r.blocks <- block
			}

			if lastBlock.Header.Number == blockRange.To {
				allBlocksDownloaded = true
			} else {
				from = lastBlock.Header.Number + 1
			}
		}
	}
}
