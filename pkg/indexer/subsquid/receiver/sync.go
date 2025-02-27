package receiver

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (r *Receiver) sync(ctx context.Context) {
	head, err := r.Head(ctx)
	if err != nil {
		return
	}
	r.MustOutput(HeadOutput).Push(head)

	if head < r.getIndexerHeight() {
		log.Warn().
			Uint64("indexer_height", r.getIndexerHeight()).
			Uint64("node_height", head).
			Msg("rollback detected by block height")
	}

	r.Log.Info().
		Uint64("indexer_block", r.getIndexerHeight()).
		Uint64("node_block", head).
		Msg("syncing...")

	startLevel := r.startLevel
	if startLevel < r.getIndexerHeight() {
		startLevel = r.getIndexerHeight()
		if r.getIndexerHeight() > 0 {
			startLevel += 1
		}
	}

	blocksToWorker, err := r.GetSqdWorkerRanges(ctx, startLevel, head)
	if err != nil {
		return
	}

	for _, blockRange := range r.SplitWorkerRanger(blocksToWorker) {
		select {
		case <-ctx.Done():
			return
		default:
			if r.checkQueue(ctx) {
				return
			}
			r.AddTask(blockRange)
		}
	}
}
