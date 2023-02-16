package indexer

import (
	"context"
	"sync"
	"time"

	starknet "github.com/dipdup-io/starknet-go-api/pkg/api"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	indexName = "dipdup_starknet_indexer"
	indexType = "rollup"
)

// Indexer -
type Indexer struct {
	cfg          Config
	outputs      map[string]*modules.Output
	queue        map[uint64]starknet.BlockWithTxs
	blocks       models.IBlock
	transactable storage.Transactable
	state        *state
	receiver     *Receiver
	wg           *sync.WaitGroup
}

// New - creates new indexer entity
func New(
	cfg Config,
	blocks models.IBlock,
	transactable storage.Transactable,
) *Indexer {
	return &Indexer{
		cfg:          cfg,
		outputs:      make(map[string]*modules.Output),
		queue:        make(map[uint64]starknet.BlockWithTxs),
		blocks:       blocks,
		transactable: transactable,
		state:        newState(nil),
		receiver:     NewReceiver(cfg),
		wg:           new(sync.WaitGroup),
	}
}

// Start -
func (indexer *Indexer) Start(ctx context.Context) {
	if err := indexer.initState(ctx); err != nil {
		log.Err(err).Msg("state initializing error")
		return
	}

	indexer.receiver.Start(ctx)

	indexer.wg.Add(1)
	go indexer.saveBlocks(ctx)

	indexer.wg.Add(1)
	go indexer.sync(ctx)
}

// Name -
func (indexer *Indexer) Name() string {
	if indexer.cfg.Name == "" {
		return indexName
	}
	return indexer.cfg.Name
}

// Close -
func (indexer *Indexer) Close() error {
	indexer.wg.Wait()
	log.Info().Msgf("closing indexer...")

	if err := indexer.receiver.Close(); err != nil {
		return err
	}

	return nil
}

func (indexer *Indexer) initState(ctx context.Context) error {
	current, err := indexer.blocks.Last(ctx)
	switch {
	case err == nil:
		indexer.state = newState(&current)
		return nil
	case indexer.blocks.IsNoRows(err):
		return nil
	default:
		return err
	}
}

func (indexer *Indexer) checkQueue(ctx context.Context) {
	for indexer.receiver.QueueSize() >= indexer.cfg.ThreadsCount {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func (indexer *Indexer) getNewBlocks(ctx context.Context) error {
	head, err := indexer.receiver.Head(ctx)
	if err != nil {
		return err
	}

	for head > indexer.state.Height() {
		log.Info().
			Uint64("indexer_block", indexer.state.Height()).
			Uint64("node_block", head).
			Msg("syncing...")

		startLevel := indexer.cfg.StartLevel
		if startLevel < indexer.state.Height()+1 {
			startLevel = indexer.state.Height() + 1
		}

		for height := startLevel; height <= head; height++ {
			select {
			case <-ctx.Done():
				return nil
			default:
				indexer.checkQueue(ctx)
				indexer.receiver.AddTask(height)
			}
		}

		for head, err = indexer.receiver.Head(ctx); err != nil; {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Err(err).Msg("receive head error")
				time.Sleep(time.Second * 5)
			}
		}
	}

	log.Info().Uint64("height", indexer.state.Height()).Msg("synced")
	return nil
}

func (indexer *Indexer) sync(ctx context.Context) {
	defer indexer.wg.Done()

	if err := indexer.getNewBlocks(ctx); err != nil {
		log.Err(err).Msg("getNewBlocks")
	}

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := indexer.getNewBlocks(ctx); err != nil {
				log.Err(err).Msg("getNewBlocks")
			}
		}
	}
}

func (indexer *Indexer) saveBlocks(ctx context.Context) {
	defer indexer.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case block := <-indexer.receiver.Results():
			indexer.queue[block.BlockNumber] = block

			next := indexer.state.Height() + 1
			if next < indexer.cfg.StartLevel+1 {
				next = indexer.cfg.StartLevel + 1
			}
			if data, ok := indexer.queue[next]; ok {
				if err := indexer.handleBlock(ctx, data); err != nil {
					log.Err(err).Msg("handle block")
				}
			}
		}
	}
}

func (indexer *Indexer) handleBlock(ctx context.Context, block starknet.BlockWithTxs) error {
	if err := indexer.handleReorg(ctx, block); err != nil {
		return err
	}

	model, err := getInternalModels(block)
	if err != nil {
		return err
	}

	if err := indexer.save(ctx, model); err != nil {
		return errors.Wrap(err, "saving block to database")
	}

	indexer.updateState(model)

	log.Info().Uint64("height", block.BlockNumber).Msg("indexed")
	delete(indexer.queue, block.BlockNumber)

	// indexer.notifyAllAboutBlock(storageData.Blocks)
	return nil
}

func (indexer *Indexer) handleReorg(ctx context.Context, block starknet.BlockWithTxs) error {
	// currentHash := indexer.state.Hash()
	// if len(currentHash) == 0 {
	// 	return nil
	// }
	// if !bytes.Equal(ethData.ParentHash().Bytes(), currentHash) {
	// 	if err := indexer.executor.Reorg(ethData.Block); err != nil {
	// 		return err
	// 	}

	// 	if err := indexer.domains.Reorg(ctx, block.Height()); err != nil {
	// 		return err
	// 	}

	// 	current, err := indexer.blocks.Last(ctx)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	indexer.state.Set(&current)

	// 	// indexer.notifyAllAboutReorg(&current)
	// }
	return nil
}

func (indexer *Indexer) updateState(block models.Block) {
	if indexer.state.Height() < block.Height {
		indexer.state.Set(block)
	}
}
