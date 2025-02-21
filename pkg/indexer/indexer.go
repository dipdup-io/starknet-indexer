package indexer

import (
	"bytes"
	"context"
	sqdAdapter "github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/adapter"
	sqdRcvr "github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver"
	"runtime"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/starknet"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/cache"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/parser/generator"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/store"
	ddConfig "github.com/dipdup-net/go-lib/config"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	indexName = "dipdup_starknet_indexer"
	indexType = "rollup"
)

// Indexer -
type Indexer struct {
	modules.BaseModule

	cfg   config.Config
	queue map[uint64]receiver.Result

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
	receiver        receiver.IReceiver
	adapter         *sqdAdapter.Adapter
	statusChecker   *statusChecker
	rollbackManager models.Rollback

	rollback      chan struct{}
	rollbackRerun chan struct{}
	rollbackWait  *sync.WaitGroup

	txWriteMutex *sync.Mutex
}

// New - creates new indexer entity
func New(
	cfg config.Config,
	storage postgres.Storage,
	datasource map[string]ddConfig.DataSource,
) (*Indexer, error) {
	indexer := &Indexer{
		BaseModule:      modules.New("indexer"),
		cfg:             cfg,
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
		rollbackManager: storage.RollbackManager,
		rollback:        make(chan struct{}, 1),
		rollbackRerun:   make(chan struct{}, 1),
		txWriteMutex:    new(sync.Mutex),
		rollbackWait:    new(sync.WaitGroup),
	}

	switch cfg.Datasource {
	case "subsquid":
		sqdReceiver, err := sqdRcvr.New(
			cfg,
			datasource,
			cfg.StartLevel,
			cfg.ThreadsCount,
			func() uint64 {
				return indexer.state.Height()
			},
		)
		if err != nil {
			return nil, err
		}

		indexer.receiver = sqdReceiver
		indexer.adapter = sqdAdapter.New(sqdReceiver.GetResults())

		if err = indexer.adapter.AttachTo(sqdReceiver, sqdRcvr.OutputName, sqdAdapter.InputName); err != nil {
			return nil, errors.Wrap(err, "while attaching adapter to receiver")
		}
	default:
		rcvr, err := receiver.NewReceiver(cfg, datasource)
		if err != nil {
			return nil, err
		}
		indexer.receiver = rcvr
	}

	indexer.CreateOutput(OutputBlocks)

	indexer.idGenerator = generator.NewIdGenerator(storage.Address, storage.Class, indexer.cache, indexer.state.Current())
	indexer.store = store.New(
		indexer.cache,
		storage.Class,
		storage.Address,
		storage.Internal,
		storage.Transfer,
		storage.Event,
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

	return indexer, nil
}

// Start -
func (indexer *Indexer) Start(ctx context.Context) {
	indexer.Log.Info().Msg("starting indexer...")
	if err := indexer.init(ctx); err != nil {
		indexer.Log.Err(err).Msg("state initializing error")
		return
	}

	indexer.receiver.Start(ctx)
	indexer.statusChecker.Start(ctx)

	switch indexer.cfg.Datasource {
	case "subsquid":
		indexer.adapter.Start(ctx)
	default:
		indexer.G.GoCtx(ctx, indexer.sync)
	}

	indexer.G.GoCtx(ctx, indexer.saveBlocks)
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
	indexer.G.Wait()
	indexer.Log.Info().Msgf("closing...")

	if err := indexer.statusChecker.Close(); err != nil {
		return err
	}

	if err := indexer.receiver.Close(); err != nil {
		return err
	}
	close(indexer.rollback)
	close(indexer.rollbackRerun)
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

func (indexer *Indexer) getNewBlocks(ctx context.Context, commonReceiver *receiver.Receiver) error {
	head, err := commonReceiver.Head(ctx)
	if err != nil {
		return err
	}

	if head < indexer.state.Height() {
		log.Warn().
			Uint64("indexer_height", indexer.state.Height()).
			Uint64("node_height", head).
			Msg("rollback detected by block height")
		if err := indexer.makeRollback(ctx, head); err != nil {
			return errors.Wrap(err, "makeRollback")
		}
	}

	for head > indexer.state.Height() {
		indexer.Log.Info().
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
			case <-indexer.rollback:
				log.Info().Msg("stop receiving blocks")
				return nil
			default:
				if indexer.checkQueue(ctx) {
					return nil
				}
				commonReceiver.AddTask(height)
			}
		}

		time.Sleep(5 * time.Second)

		for head, err = indexer.receiver.Head(ctx); err != nil; {
			select {
			case <-ctx.Done():
				return nil
			case <-indexer.rollback:
				log.Info().Msg("stop receiving blocks")
				return nil
			default:
				log.Err(err).Msg("receive head error")
				return err
			}
		}
	}

	indexer.Log.Info().Uint64("height", indexer.state.Height()).Msg("synced")
	return nil
}

func (indexer *Indexer) sync(ctx context.Context) {
	commonReceiver, ok := indexer.receiver.(*receiver.Receiver)
	if !ok {
		log.Panic().Msg("incorrect receiver type")
		return
	}
	if err := indexer.getNewBlocks(ctx, commonReceiver); err != nil {
		indexer.Log.Err(err).Msg("getNewBlocks")
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		indexer.rollbackWait.Wait()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := indexer.getNewBlocks(ctx, commonReceiver); err != nil {
				indexer.Log.Err(err).Msg("getNewBlocks")
			}
		case <-indexer.rollbackRerun:
			if err := indexer.getNewBlocks(ctx, commonReceiver); err != nil {
				indexer.Log.Err(err).Msg("getNewBlocks")
			}
		}
	}
}

func (indexer *Indexer) saveBlocks(ctx context.Context) {
	var zeroBlock bool

	for {
		select {
		case <-ctx.Done():
			return

		case result := <-indexer.receiver.Results():
			indexer.queue[result.Block.Height] = result

			if indexer.state.Height() == 0 && !zeroBlock && indexer.cfg.StartLevel == 0 {
				if data, ok := indexer.queue[0]; ok {
					if err := indexer.handleBlock(ctx, data); err != nil {
						indexer.Log.Err(err).Msg("handle block")
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
					reorgDetected, err := indexer.handleReorg(ctx, data.Block)
					if err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}
						indexer.Log.Err(err).Stack().Msg("handle reorg")
						time.Sleep(time.Second * 3)
					}

					if reorgDetected {
						break
					}

					if err := indexer.handleBlock(ctx, data); err != nil {
						if errors.Is(err, context.Canceled) {
							return
						}
						indexer.Log.Err(err).Stack().Msg("handle block")
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

func (indexer *Indexer) handleReorg(ctx context.Context, newBlock receiver.Block) (bool, error) {
	lastBlock, err := indexer.blocks.Last(ctx)
	if err != nil {
		if indexer.blocks.IsNoRows(err) {
			return false, nil
		}
		return false, err
	}

	if bytes.Equal(lastBlock.Hash, newBlock.ParentHash) {
		return false, nil
	}

	log.Warn().
		Hex("parent_hash_of_new_block", newBlock.ParentHash).
		Hex("indexer_head_block_hash", lastBlock.Hash).
		Msg("rollback detected by parent hash")

	err = indexer.makeRollback(ctx, lastBlock.Height-1)
	return true, err
}

func (indexer *Indexer) makeRollback(ctx context.Context, height uint64) error {
	for key := range indexer.queue {
		delete(indexer.queue, key)
	}

	commonReceiver, ok := indexer.receiver.(*receiver.Receiver)
	if !ok {
		return errors.Errorf("incorrect receiver type")
	}
	commonReceiver.Clear()

	if err := indexer.Rollback(ctx, height-1); err != nil {
		return errors.Wrap(err, "rollback")
	}

	if err := indexer.init(ctx); err != nil {
		return err
	}

	indexer.rollbackRerun <- struct{}{}
	return nil
}

func (indexer *Indexer) handleBlock(ctx context.Context, result receiver.Result) error {
	start := time.Now()

	parseResult, err := parser.Parse(ctx, indexer.receiver, indexer.cache, indexer.idGenerator, indexer.blocks, indexer.address, indexer.proxy, result)
	if err != nil {
		return err
	}

	var saveTime int64
	indexer.txWriteMutex.Lock()
	{
		startSave := time.Now()
		parseResult.State = indexer.updateState(parseResult.Block)
		if err := indexer.store.Save(ctx, parseResult); err != nil {
			return errors.Wrap(err, "saving block to database")
		}
		saveTime = time.Since(startSave).Milliseconds()
	}
	indexer.txWriteMutex.Unlock()

	if parseResult.Block.Status == models.StatusAcceptedOnL2 {
		indexer.statusChecker.addBlock(parseResult.Block)
	}

	delete(indexer.queue, result.Block.Height)

	l := indexer.Log.Info().
		Uint64("height", result.Block.Height).
		Int("tx_count", parseResult.Block.TxCount).
		Time("block_time", parseResult.Block.Time).
		Int64("process_time_in_ms", time.Since(start).Milliseconds()).
		Int64("save_time_in_ms", saveTime)
	if result.Block.Version != nil && *result.Block.Version != "" {
		l.Str("version", *result.Block.Version)
	}
	l.Msg("indexed")

	indexer.notifyAllAboutBlock(
		parseResult.Block,
		parseResult.Context.Addresses(),
		parseResult.Context.Tokens(),
	)
	return nil
}

func (indexer *Indexer) updateState(block models.Block) *models.State {
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
	indexer.rollbackWait.Add(1)
	defer indexer.rollbackWait.Done()

	indexer.txWriteMutex.Lock()
	defer indexer.txWriteMutex.Unlock()

	indexer.rollback <- struct{}{}

	return indexer.rollbackManager.Rollback(ctx, indexer.Name(), height)
}
