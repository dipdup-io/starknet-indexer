package sqd_receiver

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func (r *Receiver) worker(ctx context.Context, blockRange BlocksToWorker) {
	r.log.Info().
		Str("URL", blockRange.WorkerURL).
		Msg("worker handling sqd worker...")

	from := blockRange.From
	for {
		select {
		case <-ctx.Done():
			return
			// todo: indexer.rollback
		//case <-f.indexer.rollback:
		//	log.Info().Msg("stop receiving blocks")
		//	return
		default:
			blocks, err := r.api.GetBlocks(ctx, from, blockRange.WorkerURL)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				time.Sleep(time.Second)
				continue
			}

			lastBlock := blocks[len(blocks)-1]
			if lastBlock.Header.Number == blockRange.To {
				break
			}
			from = lastBlock.Header.Number + 1

			for _, block := range blocks {
				r.blocks <- block
			}

			r.log.Info().
				Uint64("From", blocks[0].Header.Number).
				Uint64("To", lastBlock.Header.Number).
				Msg("worker received blocks")
		}
	}

}
