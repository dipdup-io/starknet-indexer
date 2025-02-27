package receiver

import (
	"context"
	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	rcvr "github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/subsquid/receiver/api"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	ddConfig "github.com/dipdup-net/go-lib/config"
)

type BlocksToWorker struct {
	From      uint64
	To        uint64
	WorkerURL string
}

type GetIndexerHeight func() uint64

const (
	BlocksOutput = "blocks"
	HeadOutput   = "head_level"
	StopOutput   = "stop"
)

type Receiver struct {
	modules.BaseModule
	api              *api.Subsquid
	nodeApi          rcvr.API
	startLevel       uint64
	level            uint64
	threadsCount     int
	blocks           chan *api.SqdBlockResponse
	getIndexerHeight GetIndexerHeight
	pool             *workerpool.Pool[BlocksToWorker]
	processing       map[uint64]struct{}
	processingMx     *sync.Mutex
	result           chan rcvr.Result
	timeout          time.Duration
	mx               *sync.RWMutex
}

// New -
func New(cfg config.Config,
	ds map[string]ddConfig.DataSource,
	startLevel uint64,
	threadsCount int,
	getIndexerHeight GetIndexerHeight,
) (*Receiver, error) {
	dsCfg, ok := ds[cfg.Datasource]
	if !ok {
		return nil, errors.Errorf("unknown datasource name: %s", cfg.Datasource)
	}

	nodeCfg, ok := ds["node"]
	if !ok {
		return nil, errors.Errorf("can't access node datasource: %s", cfg.Datasource)
	}

	receiver := &Receiver{
		BaseModule:       modules.New("sqd receiver"),
		startLevel:       startLevel,
		getIndexerHeight: getIndexerHeight,
		threadsCount:     threadsCount,
		api:              api.NewSubsquid(dsCfg),
		nodeApi:          rcvr.NewNode(nodeCfg),
		blocks:           make(chan *api.SqdBlockResponse, cfg.ThreadsCount*10),
		result:           make(chan rcvr.Result, cfg.ThreadsCount*2),
		processing:       make(map[uint64]struct{}),
		processingMx:     new(sync.Mutex),
		timeout:          time.Duration(cfg.Timeout) * time.Second,
		mx:               new(sync.RWMutex),
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10 * time.Second
	}

	receiver.CreateOutput(BlocksOutput)
	receiver.CreateOutput(HeadOutput)
	receiver.CreateOutput(StopOutput)

	receiver.pool = workerpool.NewPool(receiver.worker, threadsCount)
	return receiver, nil
}

// Start -
func (r *Receiver) Start(ctx context.Context) {
	r.Log.Info().Msg("starting subsquid receiver...")
	level := r.getIndexerHeight()
	if r.startLevel > level {
		level = r.startLevel
	}

	r.setLevel(level)

	r.pool.Start(ctx)
	r.G.GoCtx(ctx, r.sync)
	r.G.GoCtx(ctx, r.sequencer)
}

// Close -
func (r *Receiver) Close() error {
	r.Log.Info().Msg("closing...")
	r.G.Wait()

	if err := r.pool.Close(); err != nil {
		return err
	}

	close(r.result)
	return nil
}

func (r *Receiver) checkQueue(ctx context.Context) bool {
	for r.pool.QueueSize() >= r.threadsCount {
		select {
		case <-ctx.Done():
			return true
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}

	return false
}

// QueueSize -
func (r *Receiver) QueueSize() int {
	return r.pool.QueueSize()
}

// AddTask -
func (r *Receiver) AddTask(blocksRange BlocksToWorker) {
	r.processingMx.Lock()
	defer r.processingMx.Unlock()

	if _, ok := r.processing[blocksRange.From]; ok {
		return
	}

	r.pool.AddTask(blocksRange)
	r.processing[blocksRange.From] = struct{}{}
}

// Results -
func (r *Receiver) Results() <-chan rcvr.Result {
	return r.result
}

// GetResults -
func (r *Receiver) GetResults() chan rcvr.Result {
	return r.result
}

// GetClass -
func (r *Receiver) GetClass(ctx context.Context, hash string) (starknetData.Class, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.nodeApi.GetClass(requestCtx, hash)
}

// Head -
func (r *Receiver) Head(ctx context.Context) (uint64, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.api.GetHead(requestCtx)
}

// GetBlockStatus -
func (r *Receiver) GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.nodeApi.GetBlockStatus(requestCtx, height)
}

func (r *Receiver) GetSqdWorkerRanges(ctx context.Context, fromLevel, height uint64) ([]BlocksToWorker, error) {
	r.Log.Info().
		Uint64("head", height).
		Msg("retrieving subsquid workers...")

	result := make([]BlocksToWorker, 0)
	currentLevel := fromLevel

	for {
		workerUrl, err := r.api.GetWorkerUrl(ctx, currentLevel)
		if err != nil {
			return nil, err
		}

		blankBlocks, err := r.api.GetBlankBlocks(ctx, currentLevel, workerUrl)
		if err != nil {
			return nil, err
		}

		lastBlock := blankBlocks[len(blankBlocks)-1]

		workerSegment := BlocksToWorker{
			From:      blankBlocks[0].Header.Number,
			To:        lastBlock.Header.Number,
			WorkerURL: workerUrl,
		}
		result = append(result, workerSegment)

		if lastBlock.Header.Number >= height {
			break
		}

		currentLevel = lastBlock.Header.Number + 1

		r.Log.Info().
			Uint64("from level", workerSegment.From).
			Uint64("to level", workerSegment.To).
			Msg("retrieved worker for blocks")
	}

	return result, nil
}

func (r *Receiver) SplitWorkerRanger(workerRanges []BlocksToWorker) []BlocksToWorker {
	var result []BlocksToWorker
	batchSize := uint64(200)

	for _, worker := range workerRanges {
		for start := worker.From; start <= worker.To; start += batchSize {
			end := start + batchSize - 1
			if end > worker.To {
				end = worker.To
			}

			result = append(result, BlocksToWorker{
				From:      start,
				To:        end,
				WorkerURL: worker.WorkerURL,
			})
		}
	}

	return result
}

func (r *Receiver) Level() uint64 {
	r.mx.RLock()
	defer r.mx.RUnlock()

	return r.level
}

func (r *Receiver) setLevel(level uint64) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.level = level
}
