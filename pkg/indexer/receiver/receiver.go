package receiver

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver/subsquid"
	"sync"
	"time"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-io/workerpool"
	ddConfig "github.com/dipdup-net/go-lib/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Result -
type Result struct {
	Block       Block
	Traces      []starknet.Trace
	StateUpdate starknetData.StateUpdate

	mx *sync.Mutex
}

func NewResult() Result {
	return Result{
		mx: new(sync.Mutex),
	}
}

func (r *Result) setBlock(block Block) {
	r.mx.Lock()
	{
		r.Block = block
	}
	r.mx.Unlock()
}

func (r *Result) setTraces(traces []starknet.Trace) {
	r.mx.Lock()
	{
		r.Traces = traces
	}
	r.mx.Unlock()
}

func (r *Result) setStateUpdates(stateUpdate starknetData.StateUpdate) {
	r.mx.Lock()
	{
		r.StateUpdate = stateUpdate
	}
	r.mx.Unlock()
}

// Receiver -
type Receiver struct {
	api          API
	fallbackAPI  API
	sqdAPI       *subsquid.Subsquid
	result       chan Result
	pool         *workerpool.Pool[uint64]
	processing   map[uint64]struct{}
	processingMx *sync.Mutex
	log          zerolog.Logger
	timeout      time.Duration
	wg           *sync.WaitGroup
}

// NewReceiver -
func NewReceiver(cfg config.Config, ds map[string]ddConfig.DataSource) (*Receiver, error) {
	dsCfg, ok := ds[cfg.Datasource]
	if !ok {
		return nil, errors.Errorf("unknown datasource name: %s", cfg.Datasource)
	}

	receiver := &Receiver{
		result:       make(chan Result, cfg.ThreadsCount*2),
		processing:   make(map[uint64]struct{}),
		processingMx: new(sync.Mutex),
		log:          log.With().Str("module", "receiver").Logger(),
		timeout:      time.Duration(cfg.Timeout) * time.Second,
		wg:           new(sync.WaitGroup),
	}

	switch cfg.Datasource {
	case "node":
		receiver.api = NewNode(dsCfg)
	case "subsquid":
		receiver.sqdAPI = subsquid.NewSubsquid(dsCfg)
	default:
		return nil, errors.Errorf("usupported datasource type: %s", cfg.Datasource)
	}

	if fallbackDs, ok := ds["fallback"]; ok && fallbackDs.URL != "" {
		receiver.fallbackAPI = NewNode(fallbackDs)
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10 * time.Second
	}

	receiver.pool = workerpool.NewPool(receiver.worker, cfg.ThreadsCount)

	return receiver, nil
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
	var (
		result = NewResult()
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go r.getBlock(ctx, blockId, &result, &wg)

	wg.Add(1)
	go r.traceBlock(ctx, blockId, &result, &wg)

	wg.Add(1)
	go r.receiveStateUpdate(ctx, blockId, &result, &wg)

	wg.Wait()

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

	if r.usingSqd() {
		return r.sqdAPI.GetHead(requestCtx)
	}

	return r.api.Head(requestCtx)
}

// GetClass -
func (r *Receiver) GetClass(ctx context.Context, hash string) (starknetData.Class, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.api.GetClass(requestCtx, hash)
}

// GetBlockStatus -
func (r *Receiver) GetBlockStatus(ctx context.Context, height uint64) (storage.Status, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.api.GetBlockStatus(requestCtx, height)
}

// Results -
func (r *Receiver) Results() <-chan Result {
	return r.result
}

// QueueSize -
func (r *Receiver) QueueSize() int {
	return r.pool.QueueSize()
}

func (r *Receiver) getStateUpdate(ctx context.Context, api API, blockId starknetData.BlockID) (starknetData.StateUpdate, error) {
	requestCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return api.GetStateUpdate(requestCtx, blockId)
}

func (r *Receiver) Clear() {
	r.pool.Clear()

	r.processingMx.Lock()
	defer r.processingMx.Unlock()

	for key := range r.processing {
		delete(r.processing, key)
	}
}

func (r *Receiver) getBlock(ctx context.Context, blockId starknetData.BlockID, result *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	api := r.api
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := api.GetBlock(ctx, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Uint64("height", *blockId.Number).Msg("get block request")
			if r.fallbackAPI != nil {
				r.log.Warn().Msg("trying fallback node...")
				api = r.fallbackAPI
			}
			time.Sleep(time.Second)
			continue
		}
		result.setBlock(response)
		break
	}
}

func (r *Receiver) usingSqd() bool {
	return r.sqdAPI != nil
}

func (r *Receiver) traceBlock(ctx context.Context, blockId starknetData.BlockID, result *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	api := r.api
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := api.TraceBlock(ctx, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Uint64("height", *blockId.Number).Msg("get block traces request")
			if r.fallbackAPI != nil {
				r.log.Warn().Msg("trying fallback node...")
				api = r.fallbackAPI
			}
			time.Sleep(time.Second)
			continue
		}
		result.setTraces(response)
		break
	}
}

func (r *Receiver) receiveStateUpdate(ctx context.Context, blockId starknetData.BlockID, result *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	api := r.api
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := r.getStateUpdate(ctx, api, blockId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.Err(err).Uint64("height", *blockId.Number).Msg("state update request")
			if r.fallbackAPI != nil {
				r.log.Warn().Msg("trying fallback node...")
				api = r.fallbackAPI
			}
			time.Sleep(time.Second)
			continue
		}
		result.setStateUpdates(response)
		break
	}
}

func (r *Receiver) GetSqdData(ctx context.Context, startLevel uint64) {
	r.sqdAPI.GetData(ctx, startLevel)
	//for {
	//	select {
	//	case <-ctx.Done():
	//		return
	//	default:
	//	}
	//
	//	response, err := api.GetData(ctx, blockId)
	//	if err != nil {
	//		if errors.Is(err, context.Canceled) {
	//			return
	//		}
	//		r.log.Err(err).Uint64("height", *blockId.Number).Msg("get block request")
	//		if r.fallbackAPI != nil {
	//			r.log.Warn().Msg("trying fallback node...")
	//			api = r.fallbackAPI
	//		}
	//		time.Sleep(time.Second)
	//		continue
	//	}
	//	result.setBlock(response)
	//	break
	//}
}
