package sqd_receiver

import (
	"context"
	rcvr "github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/sqd_receiver/api"
	"github.com/dipdup-io/workerpool"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	OutputName = "blocks"
	StopOutput = "stop"
)

type Receiver struct {
	modules.BaseModule
	api              *api.Subsquid
	startLevel       uint64
	level            uint64
	threadsCount     int
	blocks           chan *api.SqdBlockResponse
	getIndexerHeight GetIndexerHeight
	pool             *workerpool.Pool[BlocksToWorker]
	processing       map[uint64]struct{}
	processingMx     *sync.Mutex
	result           chan rcvr.Result
	log              zerolog.Logger
	timeout          time.Duration
	wg               *sync.WaitGroup
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

	receiver := &Receiver{
		BaseModule:       modules.New("subsquid receiver"),
		startLevel:       startLevel,
		getIndexerHeight: getIndexerHeight,
		threadsCount:     threadsCount,
		api:              api.NewSubsquid(dsCfg),
		blocks:           make(chan *api.SqdBlockResponse, cfg.ThreadsCount*10),
		processing:       make(map[uint64]struct{}),
		processingMx:     new(sync.Mutex),
		log:              log.With().Str("module", "subsquid_receiver").Logger(),
		timeout:          time.Duration(cfg.Timeout) * time.Second,
		wg:               new(sync.WaitGroup),
		mx:               new(sync.RWMutex),
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10 * time.Second
	}

	receiver.CreateOutput(OutputName)
	receiver.CreateOutput(StopOutput)

	receiver.pool = workerpool.NewPool(receiver.worker, cfg.ThreadsCount)
	return receiver, nil
}

// Start -
func (r *Receiver) Start(ctx context.Context) {
	r.log.Info().Msg("starting subsquid receiver...")
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
	r.log.Info().Msg("closing...")
	r.wg.Wait()

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
			time.Sleep(time.Millisecond * 100)
		}
	}

	return false
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

func (r *Receiver) GetSqdWorkerRanges(ctx context.Context, fromLevel, height uint64) ([]BlocksToWorker, error) {
	r.log.Info().
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

		if lastBlock.Header.Number == height {
			break
		}

		currentLevel = lastBlock.Header.Number + 1

		r.log.Info().
			Uint64("from level", workerSegment.From).
			Uint64("to level", workerSegment.To).
			Msg("retrieved worker for blocks")
	}

	return result, nil
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
