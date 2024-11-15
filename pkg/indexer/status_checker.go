package indexer

import (
	"context"
	"time"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/workerpool"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type acceptedOnL2 struct {
	Height uint64
}

type statusChecker struct {
	acceptedOnL2   queue[acceptedOnL2]
	blocks         storage.IBlock
	declares       storage.IDeclare
	deploys        storage.IDeploy
	deployAccounts storage.IDeployAccount
	invoke         storage.IInvoke
	l1Handlers     storage.IL1Handler
	transactable   sdk.Transactable
	receiver       *receiver.Receiver
	log            zerolog.Logger
	g              workerpool.Group
}

func newStatusChecker(
	receiver *receiver.Receiver,
	blocks storage.IBlock,
	declares storage.IDeclare,
	deploys storage.IDeploy,
	deployAccounts storage.IDeployAccount,
	invoke storage.IInvoke,
	l1Handlers storage.IL1Handler,
	transactable sdk.Transactable,
) *statusChecker {
	return &statusChecker{
		acceptedOnL2:   newQueue[acceptedOnL2](),
		receiver:       receiver,
		blocks:         blocks,
		declares:       declares,
		deploys:        deploys,
		deployAccounts: deployAccounts,
		invoke:         invoke,
		l1Handlers:     l1Handlers,
		transactable:   transactable,
		log:            log.With().Str("module", "status_checker").Logger(),
		g:              workerpool.NewGroup(),
	}
}

// Start -
func (checker *statusChecker) Start(ctx context.Context) {
	checker.g.GoCtx(ctx, checker.start)
}

func (checker *statusChecker) start(ctx context.Context) {
	if err := checker.init(ctx); err != nil {
		checker.log.Err(err).Msg("checker init")
		return
	}

	if err := checker.check(ctx); err != nil {
		checker.log.Err(err).Msg("check block status")
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := checker.check(ctx); err != nil {
				checker.log.Err(err).Msg("check block status")
			}
		}
	}
}

// Close -
func (checker *statusChecker) Close() error {
	checker.g.Wait()
	return nil
}

func (checker *statusChecker) init(ctx context.Context) error {
	var (
		offset uint64
		end    bool
		limit  = uint64(50)
	)

	for !end {
		blocks, err := checker.blocks.ByStatus(ctx, storage.StatusAcceptedOnL2, limit, offset, sdk.SortOrderAsc)
		if err != nil {
			return err
		}
		count := uint64(len(blocks))
		offset += count
		end = limit != count

		for i := range blocks {
			checker.acceptedOnL2.Push(acceptedOnL2{
				Height: blocks[i].Height,
			})
		}
	}

	return nil
}

func (checker *statusChecker) check(ctx context.Context) error {
	for {
		if checker.acceptedOnL2.IsEmpty() {
			return nil
		}

		item, err := checker.acceptedOnL2.First()
		if err != nil {
			return err
		}

		status, err := checker.getStatus(ctx, item)
		if err != nil {
			return err
		}

		if status != storage.StatusAcceptedOnL1 {
			return nil
		}

		if err := checker.update(ctx, item.Height, status); err != nil {
			return err
		}

		if _, err := checker.acceptedOnL2.Pop(); err != nil {
			return err
		}

		checker.log.Info().Str("status", status.String()).Uint64("height", item.Height).Msg("update block status")
	}
}

func (checker *statusChecker) addBlock(block storage.Block) {
	checker.acceptedOnL2.Push(acceptedOnL2{block.Height})
}

func (checker *statusChecker) getStatus(ctx context.Context, item acceptedOnL2) (storage.Status, error) {
	return checker.receiver.GetBlockStatus(ctx, item.Height)
}

func (checker *statusChecker) update(ctx context.Context, height uint64, status storage.Status) error {
	tx, err := postgres.BeginTransaction(ctx, checker.transactable)
	if err != nil {
		return err
	}

	for _, model := range []sdk.Model{
		&storage.Block{},
		&storage.Declare{},
		&storage.DeployAccount{},
		&storage.Deploy{},
		&storage.Invoke{},
		&storage.L1Handler{},
		&storage.Internal{},
		&storage.Fee{},
	} {
		if err := tx.UpdateStatus(ctx, height, status, model); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	return tx.Flush(ctx)
}
