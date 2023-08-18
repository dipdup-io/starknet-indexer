package receiver

import (
	"context"
	"errors"
	"sync"
	"time"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	starknetRpc "github.com/dipdup-io/starknet-go-api/pkg/rpc"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-io/workerpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Result -
type Result struct {
	Block       starknet.Block
	Trace       starknet.TraceResponse
	StateUpdate starknetData.StateUpdate
}

// Receiver -
type Receiver struct {
	api          starknet.API
	rpc          *starknetRpc.API
	result       chan Result
	pool         *workerpool.Pool[uint64]
	processing   map[uint64]struct{}
	processingMx *sync.Mutex
	log          zerolog.Logger
	timeout      time.Duration
	wg           *sync.WaitGroup
}

// NewReceiver -
func NewReceiver(cfg config.Config) *Receiver {
	opts := make([]starknet.ApiOption, 0)
	if cfg.Sequencer.Rps > 0 {
		opts = append(opts, starknet.WithRateLimit(cfg.Sequencer.Rps))
	}
	if cfg.CacheDir != "" {
		opts = append(opts, starknet.WithCacheInFS(cfg.CacheDir))
	}

	api := starknet.NewAPI(cfg.Sequencer.Gateway, cfg.Sequencer.FeederGateway, opts...)

	receiver := &Receiver{
		api:          api,
		result:       make(chan Result, cfg.ThreadsCount*2),
		processing:   make(map[uint64]struct{}),
		processingMx: new(sync.Mutex),
		log:          log.With().Str("module", "receiver").Logger(),
		timeout:      time.Duration(cfg.Timeout) * time.Second,
		wg:           new(sync.WaitGroup),
	}

	if cfg.Node != nil && cfg.Node.Url != "" {
		rpc := starknetRpc.NewAPI(cfg.Node.Url, starknetRpc.WithRateLimit(cfg.Node.Rps))
		receiver.rpc = &rpc
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10 * time.Second
	}

	receiver.pool = workerpool.NewPool(receiver.worker, cfg.ThreadsCount)

	return receiver
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

// Start -
func (r *Receiver) Start(ctx context.Context) {
	r.pool.Start(ctx)
}

func (r *Receiver) worker(ctx context.Context, height uint64) {
	start := time.Now()
	blockId := starknetData.BlockID{
		Number: &height,
	}
	var result Result
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := r.api.GetBlock(ctx, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Msg("get block request")
			time.Sleep(time.Second)
			continue
		}
		result.Block = response
		break
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := r.api.TraceBlock(ctx, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Msg("get block traces request")
			time.Sleep(time.Second)
			continue
		}
		result.Trace = response
		break
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := r.getStateUpdate(ctx, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Msg("state update request")
			time.Sleep(time.Second)
			continue
		}
		result.StateUpdate = response
		break
	}

	r.log.Info().Uint64("height", height).Int64("ms", time.Since(start).Milliseconds()).Msg("received block data")
	r.result <- result
	r.processingMx.Lock()
	{
		delete(r.processing, height)
	}
	r.processingMx.Unlock()
}

// AddTask -
func (r *Receiver) AddTask(height uint64) {
	r.processingMx.Lock()
	defer r.processingMx.Unlock()

	if _, ok := r.processing[height]; ok {
		return
	}

	r.pool.AddTask(height)
	r.processing[height] = struct{}{}
}

// Head -
func (r *Receiver) Head(ctx context.Context) (uint64, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if r.rpc != nil {
		response, err := r.rpc.BlockNumber(requestCtx)
		if err != nil {
			return 0, err
		}
		return response.Result, nil
	}

	response, err := r.api.GetBlock(requestCtx, starknetData.BlockID{
		String: starknetData.Latest,
	})
	if err != nil {
		return 0, err
	}
	return response.BlockNumber, nil
}

// GetClass -
func (r *Receiver) GetClass(ctx context.Context, hash string) (starknetData.Class, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	blockId := starknetData.BlockID{
		String: starknetData.Latest,
	}

	return r.api.GetClassByHash(requestCtx, blockId, hash)
}

// TransactionStatus -
func (r *Receiver) TransactionStatus(ctx context.Context, hash string) (storage.Status, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	response, err := r.api.GetTransactionStatus(requestCtx, hash)
	if err != nil {
		return storage.StatusUnknown, err
	}

	return storage.NewStatus(response.Status), nil
}

// GetBlockStatus -
func (r *Receiver) GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	blockId := starknetData.BlockID{
		Number: &height,
	}

	if r.rpc != nil {
		response, err := r.rpc.GetBlockWithTxHashes(requestCtx, blockId)
		if err != nil {
			return storage.StatusUnknown, err
		}
		return storage.NewStatus(response.Result.Status), nil
	}

	response, err := r.api.GetBlock(requestCtx, blockId)
	if err != nil {
		return storage.StatusUnknown, err
	}
	return storage.NewStatus(response.Status), nil
}

// Results -
func (r *Receiver) Results() <-chan Result {
	return r.result
}

// QueueSize -
func (r *Receiver) QueueSize() int {
	return r.pool.QueueSize()
}

func (r *Receiver) getStateUpdate(ctx context.Context, blockId starknetData.BlockID) (starknetData.StateUpdate, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if r.rpc != nil {
		response, err := r.rpc.GetStateUpdate(requestCtx, blockId)
		if err != nil {
			return starknetData.StateUpdate{}, err
		}
		return response.Result.ToStateUpdate(), nil
	}

	return r.api.GetStateUpdate(requestCtx, blockId)
}

func (r *Receiver) Clear() {
	r.pool.Clear()

	r.processingMx.Lock()
	defer r.processingMx.Unlock()

	for key := range r.processing {
		delete(r.processing, key)
	}
}
