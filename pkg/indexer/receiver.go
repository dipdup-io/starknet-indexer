package indexer

import (
	"context"
	"errors"
	"sync"
	"time"

	starknet "github.com/dipdup-io/starknet-go-api/pkg/api"
	"github.com/dipdup-net/workerpool"
	"github.com/rs/zerolog/log"
)

// Receiver -
type Receiver struct {
	api     starknet.API
	result  chan starknet.BlockWithTxs
	pool    *workerpool.Pool[uint64]
	timeout uint64
	wg      *sync.WaitGroup
}

// NewReceiver -
func NewReceiver(cfg Config) *Receiver {
	receiver := &Receiver{
		api:     starknet.NewAPI(cfg.BaseURL),
		result:  make(chan starknet.BlockWithTxs, cfg.ThreadsCount),
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
	for {
		response, err := r.api.GetBlockWithTxs(ctx, starknet.BlockFilter{
			Number: height,
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Err(err).Msg("block request")
			time.Sleep(time.Second)
			continue
		}
		r.result <- response.Result
		return
	}
}

// AddTask -
func (r *Receiver) AddTask(height uint64) {
	r.pool.AddTask(height)
}

// Head -
func (r *Receiver) Head(ctx context.Context) (uint64, error) {
	response, err := r.api.BlockNumber(ctx, starknet.WithTimeout(r.timeout))
	if err != nil {
		return 0, err
	}
	return response.Result, nil
}

// Results -
func (r *Receiver) Results() <-chan starknet.BlockWithTxs {
	return r.result
}

// QueueSize -
func (r *Receiver) QueueSize() int {
	return r.pool.QueueSize()
}
