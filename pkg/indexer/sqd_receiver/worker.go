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

	// todo: move to config
	var batchSize uint64 = 1000
	from := blockRange.From
	to := blockRange.From + batchSize

	for {
		select {
		case <-ctx.Done():
			return
			// todo: indexer.rollback
		//case <-f.indexer.rollback:
		//	log.Info().Msg("stop receiving blocks")
		//	return
		default:
			blocks, err := r.api.GetBlocks(ctx, from, to, blockRange.WorkerURL)
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
			to = from + batchSize
			if to > blockRange.To {
				to = blockRange.To
			}
			r.log.Info().
				Uint64("From", blocks[0].Header.Number).
				Uint64("To", lastBlock.Header.Number).
				Msg("worker received blocks")
		}
	}

}
