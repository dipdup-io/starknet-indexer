package indexer

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/generator"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/store"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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
	invoke         models.IInvoke
	l1Handlers     models.IL1Handler
	storageDiffs   models.IStorageDiff
	proxy          models.IProxy
	stateRepo      models.IState
	transactable   sdk.Transactable
	store          *store.Store
	cache          *cache.Cache

	state           *state
	idGenerator     *generator.IdGenerator
	receiver        *receiver.Receiver
	statusChecker   *statusChecker
	rollbackManager models.Rollback

	log zerolog.Logger

	txWriteMutex *sync.Mutex
	wg           *sync.WaitGroup
}

// New - creates new indexer entity
func New(
	cfg config.Config,
	storage postgres.Storage,
) *Indexer {
	indexer := &Indexer{
		cfg: cfg,
		outputs: map[string]*modules.Output{
			OutputBlocks: modules.NewOutput(OutputBlocks),
		},
		queue:           make(map[uint64]receiver.Result),
		stateRepo:       storage.State,
		address:         storage.Address,
		blocks:          storage.Blocks,
		declares:        storage.Declare,
		deploys:         storage.Deploy,
		deployAccounts:  storage.DeployAccount,
		invoke:          storage.Invoke,
		l1Handlers:      storage.L1Handler,
		classes:         storage.Class,
		storageDiffs:    storage.StorageDiff,
		transactable:    storage.Transactable,
		proxy:           storage.Proxy,
		state:           newState(nil),
		cache:           cache.New(storage.Address, storage.Class, storage.Proxy),
		receiver:        receiver.NewReceiver(cfg),
		rollbackManager: storage.RollbackManager,
		log:             log.With().Str("module", "indexer").Logger(),
		txWriteMutex:    new(sync.Mutex),
		wg:              new(sync.WaitGroup),
	}

	indexer.idGenerator = generator.NewIdGenerator(storage.Address, storage.Class, indexer.cache, indexer.state.Current())
	indexer.store = store.New(
		indexer.cache,
		storage.Class,
		storage.Address,
		storage.Internal,
		storage.Transfer,
		storage.Event,
		storage.StorageDiff,
		storage.Invoke,
		storage.Deploy,
		storage.L1Handler,
		storage.Fee,
		storage.Transactable,
		storage.PartitionManager)

	indexer.statusChecker = newStatusChecker(
		indexer.receiver,
		storage.Blocks,
		storage.Declare,
		storage.Deploy,
		storage.DeployAccount,
		storage.Invoke,
		storage.L1Handler,
		storage.Transactable,
	)

	return indexer
}

// Start -
func (indexer *Indexer) Start(ctx context.Context) {
	indexer.log.Info().Msg("starting indexer...")
	if err := indexer.init(ctx); err != nil {
		indexer.log.Err(err).Msg("state initializing error")
		return
	}

	indexer.receiver.Start(ctx)

	indexer.statusChecker.Start(ctx)

	if err := indexer.fixClassAbi(ctx); err != nil {
		indexer.log.Err(err).Msg("recovering class abi error")
		return
	}

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
	indexer.log.Info().Msgf("closing...")

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

	state, err := indexer.stateRepo.ByName(ctx, indexer.Name())
	switch {
	case err == nil:
		indexer.state.Set(state)
		indexer.idGenerator.Init()
		return nil
	case indexer.blocks.IsNoRows(err):
		state := indexer.state.Current()
		state.Name = indexer.Name()
		return indexer.stateRepo.Save(ctx, state)
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
		indexer.log.Info().
			Uint64("indexer_block", indexer.state.Height()).
			Uint64("node_block", head).
			Msg("syncing...")

		startLevel := indexer.cfg.StartLevel
		if startLevel < indexer.state.Height() {
			startLevel = indexer.state.Height()
			if indexer.state.Height() > 0 {
				startLevel += 1
			}
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
				time.Sleep(time.Second * 30)
			}
		}
	}

	indexer.log.Info().Uint64("height", indexer.state.Height()).Msg("synced")
	return nil
}

func (indexer *Indexer) sync(ctx context.Context) {
	defer indexer.wg.Done()

	if err := indexer.getNewBlocks(ctx); err != nil {
		indexer.log.Err(err).Msg("getNewBlocks")
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := indexer.getNewBlocks(ctx); err != nil {
				indexer.log.Err(err).Msg("getNewBlocks")
			}
		}
	}
}

func (indexer *Indexer) saveBlocks(ctx context.Context) {
	defer indexer.wg.Done()

	var zeroBlock bool

	for {
		select {
		case <-ctx.Done():
			return

		case result := <-indexer.receiver.Results():
			indexer.queue[result.Block.BlockNumber] = result

			if indexer.state.Height() == 0 && !zeroBlock {
				if data, ok := indexer.queue[0]; ok {
					if err := indexer.handleBlock(ctx, data); err != nil {
						indexer.log.Err(err).Msg("handle block")
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

			for {
				if data, ok := indexer.queue[next]; ok {
					if err := indexer.handleBlock(ctx, data); err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}
						indexer.log.Err(err).Stack().Msg("handle block")
						time.Sleep(time.Second * 3)
					}
					if next%25 == 0 {
						runtime.GC()
					}

					next = indexer.state.Height() + 1
				} else {
					break
				}
			}
		}
	}
}

func (indexer *Indexer) handleBlock(ctx context.Context, result receiver.Result) error {
	start := time.Now()

	parseResult, err := parser.Parse(ctx, indexer.receiver, indexer.cache, indexer.idGenerator, indexer.blocks, indexer.proxy, result)
	if err != nil {
		return err
	}

	var saveTime int64
	indexer.txWriteMutex.Lock()
	{
		startSave := time.Now()
		parseResult.State = indexer.updateState(ctx, parseResult.Block, len(parseResult.Context.Classes()))
		if err := indexer.store.Save(ctx, parseResult); err != nil {
			return errors.Wrap(err, "saving block to database")
		}
		saveTime = time.Since(startSave).Milliseconds()
	}
	indexer.txWriteMutex.Unlock()

	if parseResult.Block.Status == models.StatusAcceptedOnL2 {
		indexer.statusChecker.addBlock(parseResult.Block)
	}

	delete(indexer.queue, result.Block.BlockNumber)

	l := indexer.log.Info().
		Uint64("height", result.Block.BlockNumber).
		Int("tx_count", parseResult.Block.TxCount).
		Time("block_time", parseResult.Block.Time).
		Int64("process_time_in_ms", time.Since(start).Milliseconds()).
		Int64("save_time_in_ms", saveTime)
	if result.Block.StarknetVersion != nil {
		l.Str("version", *result.Block.StarknetVersion)
	}
	l.Msg("indexed")

	indexer.notifyAllAboutBlock(parseResult.Block, parseResult.Context.Addresses())
	return nil
}

func (indexer *Indexer) updateState(ctx context.Context, block models.Block, classesCount int) *models.State {
	state := indexer.state.Current()
	if indexer.state.Height() < block.Height {
		state.LastHeight = block.Height
		state.LastTime = block.Time
		state.DeclaresCount += uint64(block.DeclareCount)
		state.DeploysCount += uint64(block.DeployCount)
		state.DeployAccountsCount += uint64(block.DeployAccountCount)
		state.InvokesCount += uint64(block.InvokeCount)
		state.L1HandlersCount += uint64(block.L1HandlerCount)
		state.TxCount += uint64(block.TxCount)
	}
	return state
}

// Rollback -
func (indexer *Indexer) Rollback(ctx context.Context, height uint64) error {
	indexer.txWriteMutex.Lock()
	defer indexer.txWriteMutex.Unlock()

	indexer.log.Info().Uint64("new_height", height).Msg("rollback starting...")
	return indexer.rollbackManager.Rollback(ctx, indexer.Name(), height)
}

func (indexer *Indexer) fixClassAbi(ctx context.Context) error {
	tx, err := indexer.transactable.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	classes, err := indexer.classes.GetUnresolved(ctx)
	if err != nil {
		return tx.HandleError(ctx, err)
	}

	log.Info().Msgf("found %d unresolved classes", len(classes))

	for i := range classes {
		hash := data.NewFeltFromBytes(classes[i].Hash)
		class, err := indexer.receiver.GetClass(ctx, hash.String())
		if err != nil {
			return tx.HandleError(ctx, err)
		}
		classes[i].Abi = models.Bytes(class.RawAbi)
		if err := tx.Update(ctx, &classes[i]); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	if err := tx.Flush(ctx); err != nil {
		return tx.HandleError(ctx, err)
	}
	return nil
}
