package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/store"
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
	cfg     config.Config
	outputs map[string]*modules.Output
	queue   map[uint64]receiver.Result

	address        models.IAddress
	blocks         models.IBlock
	classes        models.IClass
	declares       models.IDeclare
	deploys        models.IDeploy
	deployAccounts models.IDeployAccount
	invokeV0       models.IInvokeV0
	invokeV1       models.IInvokeV1
	l1Handlers     models.IL1Handler
	storageDiffs   models.IStorageDiff

	store *store.Store
	cache *cache.Cache

	state         *state
	idGenerator   *parser.IdGenerator
	receiver      *receiver.Receiver
	statusChecker *statusChecker
	wg            *sync.WaitGroup
}

// New - creates new indexer entity
func New(
	cfg config.Config,
	address models.IAddress,
	blocks models.IBlock,
	declares models.IDeclare,
	deploys models.IDeploy,
	deployAccounts models.IDeployAccount,
	invokeV0 models.IInvokeV0,
	invokeV1 models.IInvokeV1,
	l1Handlers models.IL1Handler,
	classes models.IClass,
	storageDiffs models.IStorageDiff,
	transactable storage.Transactable,
) *Indexer {
	indexer := &Indexer{
		cfg:            cfg,
		outputs:        make(map[string]*modules.Output),
		queue:          make(map[uint64]receiver.Result),
		address:        address,
		blocks:         blocks,
		declares:       declares,
		deploys:        deploys,
		deployAccounts: deployAccounts,
		invokeV0:       invokeV0,
		invokeV1:       invokeV1,
		l1Handlers:     l1Handlers,
		classes:        classes,
		storageDiffs:   storageDiffs,
		state:          newState(nil),
		cache:          cache.New(address, classes),
		receiver:       receiver.NewReceiver(cfg),
		wg:             new(sync.WaitGroup),
	}

	indexer.idGenerator = parser.NewIdGenerator(address, classes, indexer.cache)
	indexer.store = store.New(indexer.cache, classes, address, transactable)
	indexer.statusChecker = newStatusChecker(
		indexer.receiver,
		blocks,
		declares,
		deploys,
		deployAccounts,
		invokeV0,
		invokeV1,
		l1Handlers,
		transactable,
	)

	return indexer
}

// Start -
func (indexer *Indexer) Start(ctx context.Context) {
	if err := indexer.init(ctx); err != nil {
		log.Err(err).Msg("state initializing error")
		return
	}

	indexer.receiver.Start(ctx)

	indexer.statusChecker.Start(ctx)

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

	if err := indexer.statusChecker.Close(); err != nil {
		return err
	}

	if err := indexer.receiver.Close(); err != nil {
		return err
	}

	return nil
}

func (indexer *Indexer) init(ctx context.Context) error {
	if _, err := starknet.Interfaces(indexer.cfg.ClassInterfacesDir); err != nil {
		return err
	}

	current, err := indexer.blocks.Last(ctx)
	switch {
	case err == nil:
		indexer.state = newState(&current)
		if err := indexer.idGenerator.Init(ctx); err != nil {
			return err
		}

		return nil
	case indexer.blocks.IsNoRows(err):
		return nil
	default:
		return err
	}
}

func (indexer *Indexer) checkQueue(ctx context.Context) bool {
	for indexer.receiver.QueueSize() >= indexer.cfg.ThreadsCount {
		select {
		case <-ctx.Done():
			return true
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

	return false
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
		if startLevel < indexer.state.Height() {
			startLevel = indexer.state.Height()
		}

		for height := startLevel; height <= head; height++ {
			select {
			case <-ctx.Done():
				return nil
			default:
				if indexer.checkQueue(ctx) {
					return nil
				}
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

	ticker := time.NewTicker(time.Microsecond * 100)
	defer ticker.Stop()

	var zeroBlock bool

	for {
		select {
		case <-ctx.Done():
			return

		case result := <-indexer.receiver.Results():
			indexer.queue[result.Block.BlockNumber] = result

		case <-ticker.C:
			if indexer.state.Height() == 0 && !zeroBlock {
				if data, ok := indexer.queue[0]; ok {
					if err := indexer.handleBlock(ctx, data); err != nil {
						log.Err(err).Msg("handle block")
					}
					zeroBlock = true
				} else {
					continue
				}
			}

			next := indexer.state.Height() + 1
			if next < indexer.cfg.StartLevel+1 {
				next = indexer.cfg.StartLevel + 1
			}
			if data, ok := indexer.queue[next]; ok {
				if err := indexer.handleBlock(ctx, data); err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					log.Err(err).Msg("handle block")
					time.Sleep(time.Second * 10)
				}
			}
		}
	}
}

func (indexer *Indexer) handleBlock(ctx context.Context, result receiver.Result) error {
	parser := parser.New(indexer.receiver, indexer.cache, indexer.idGenerator, indexer.storageDiffs)
	parseResult, err := parser.Parse(ctx, result)
	if err != nil {
		return err
	}

	if err := indexer.store.Save(ctx, parseResult); err != nil {
		return errors.Wrap(err, "saving block to database")
	}

	if parseResult.Block.Status == models.StatusAcceptedOnL2 {
		indexer.statusChecker.addBlock(parseResult.Block)
	}

	indexer.updateState(parseResult.Block)

	log.Info().Uint64("height", result.Block.BlockNumber).Msg("indexed")
	delete(indexer.queue, result.Block.BlockNumber)

	// indexer.notifyAllAboutBlock(storageData.Blocks)
	return nil
}

func (indexer *Indexer) updateState(block models.Block) {
	if indexer.state.Height() < block.Height {
		indexer.state.Set(block)
	}
}
