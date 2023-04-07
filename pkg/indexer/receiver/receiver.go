package receiver

import (
	"context"
	"errors"
	"sync"
	"time"

	starknetData "github.com/dipdup-io/starknet-go-api/pkg/data"
	"github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	starknet "github.com/dipdup-io/starknet-go-api/pkg/sequencer"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/config"
	"github.com/dipdup-net/workerpool"
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
	api     starknet.API
	result  chan Result
	pool    *workerpool.Pool[uint64]
	timeout uint64
	wg      *sync.WaitGroup
}

// NewReceiver -
func NewReceiver(cfg config.Config) *Receiver {
	opts := make([]starknet.ApiOption, 0)
	if cfg.RequestsPerSecond > 0 {
		opts = append(opts, starknet.WithRateLimit(cfg.RequestsPerSecond))
	}
	if cfg.CacheDir != "" {
		opts = append(opts, starknet.WithCacheInFS(cfg.CacheDir))
	}

	api := starknet.NewAPI(cfg.Gateway, cfg.FeederGateway, opts...)

	receiver := &Receiver{
		api:     api,
		result:  make(chan Result, cfg.ThreadsCount*2),
		timeout: cfg.Timeout,
		wg:      new(sync.WaitGroup),
	}

	if receiver.timeout == 0 {
		receiver.timeout = 10
	}

	receiver.pool = workerpool.NewPool(receiver.worker, cfg.ThreadsCount)

	return receiver
}

// Close -
func (r *Receiver) Close() error {
	log.Info().Msg("closing receiver...")
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
	var result Result
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		response, err := r.api.GetBlock(ctx, starknetData.BlockID{
			Number: &height,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Err(err).Msg("get block request")
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

		response, err := r.api.TraceBlock(ctx, starknetData.BlockID{
			Number: &height,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Err(err).Msg("get block traces request")
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

		requestCtx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		response, err := r.api.GetStateUpdate(requestCtx, starknetData.BlockID{
			Number: &height,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Err(err).Msg("state update request")
			time.Sleep(time.Second)
			continue
		}
		result.StateUpdate = response
		break
	}

	r.result <- result
}

// AddTask -
func (r *Receiver) AddTask(height uint64) {
	r.pool.AddTask(height)
}

// Head -
func (r *Receiver) Head(ctx context.Context) (uint64, error) {
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeout)*time.Second)
	defer cancel()

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
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeout)*time.Second)
	defer cancel()

	response, err := r.api.GetClassByHash(requestCtx, starknetData.BlockID{
		String: starknetData.Latest,
	}, hash)
	if err != nil {
		return starknetData.Class{}, err
	}
	return response, nil
}

// TransactionStatus -
func (r *Receiver) TransactionStatus(ctx context.Context, hash string) (storage.Status, error) {
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeout)*time.Second)
	defer cancel()

	response, err := r.api.GetTransactionStatus(requestCtx, hash)
	if err != nil {
		return storage.StatusUnknown, err
	}

	return storage.NewStatus(response.Status), nil
}

// GetBlock -
func (r *Receiver) GetBlock(ctx context.Context, height uint64) (sequencer.Block, error) {
	requestCtx, cancel := context.WithTimeout(ctx, time.Duration(r.timeout)*time.Second)
	defer cancel()

	return r.api.GetBlock(requestCtx, starknetData.BlockID{
		Number: &height,
	})
}

// Results -
func (r *Receiver) Results() <-chan Result {
	return r.result
}

// QueueSize -
func (r *Receiver) QueueSize() int {
	return r.pool.QueueSize()
}
