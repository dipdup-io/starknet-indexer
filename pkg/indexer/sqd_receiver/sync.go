package sqd_receiver

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (r *Receiver) sync(ctx context.Context) {
	head, err := r.api.GetHead(ctx)
	if err != nil {
		return
	}

	if head < r.getIndexerHeight() {
		log.Warn().
			Uint64("indexer_height", r.getIndexerHeight()).
			Uint64("node_height", head).
			Msg("rollback detected by block height")

		// todo: makeRollback
		//if err := f.indexer.makeRollback(ctx, head); err != nil {
		//	return errors.Wrap(err, "makeRollback")
		//}
	}

	r.log.Info().
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

	for _, blockRange := range blocksToWorker {
		select {
		case <-ctx.Done():
			return
		// todo: f.indexer.rollback
		//case <-f.indexer.rollback:
		//	log.Info().Msg("stop receiving blocks")
		//	return nil
		default:
			if r.checkQueue(ctx) {
				return
			}
			r.AddTask(blockRange)
		}
	}

	r.log.Info().Uint64("height", r.getIndexerHeight()).Msg("synced")
	return
}
