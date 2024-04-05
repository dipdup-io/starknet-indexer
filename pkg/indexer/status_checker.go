package indexer

import (
	"context"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-io/workerpool"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	errCantFindTransactionsInBlock = errors.New("ca't find transactions in block")
)

type acceptedOnL2 struct {
	Height          uint64
	TransactionHash []byte
}

func newAcceptedOnL2FromIndexingBlock(block storage.Block) acceptedOnL2 {
	a := acceptedOnL2{
		Height: block.Height,
	}

	if block.InvokeCount > 0 {
		a.TransactionHash = block.Invoke[0].Hash
		return a
	}

	if block.DeclareCount > 0 {
		a.TransactionHash = block.Declare[0].Hash
		return a
	}

	if block.DeployCount > 0 {
		a.TransactionHash = block.Deploy[0].Hash
		return a
	}

	if block.DeployAccountCount > 0 {
		a.TransactionHash = block.DeployAccount[0].Hash
		return a
	}

	if block.L1HandlerCount > 0 {
		a.TransactionHash = block.L1Handler[0].Hash
		return a
	}

	return a
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
			if err := checker.addIndexedBlockToQueue(ctx, blocks[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func byHeight[T sdk.Model, F any](ctx context.Context, src storage.Filterable[T, F], fltr F) (t T, err error) {
	tx, err := src.Filter(ctx, []F{fltr}, storage.WithLimitFilter(1), storage.WithMultiSort("time desc", "id desc"))
	if err != nil {
		return t, err
	}
	if len(tx) == 0 {
		return t, errors.Wrapf(
			errCantFindTransactionsInBlock,
			"model = %s", t.TableName())
	}

	return tx[0], nil
}

func (checker *statusChecker) addIndexedBlockToQueue(ctx context.Context, block storage.Block) error {
	if block.InvokeCount > 0 {
		tx, err := byHeight[storage.Invoke, storage.InvokeFilter](ctx, checker.invoke, storage.InvokeFilter{
			Height: storage.IntegerFilter{
				Eq: block.Height,
			},
		})
		if err != nil {
			return err
		}
		checker.acceptedOnL2.Push(acceptedOnL2{
			TransactionHash: tx.Hash,
			Height:          tx.Height,
		})
		return nil
	}
	if block.DeployCount > 0 {
		tx, err := byHeight[storage.Deploy, storage.DeployFilter](ctx, checker.deploys, storage.DeployFilter{
			Height: storage.IntegerFilter{
				Eq: block.Height,
			},
		})
		if err != nil {
			return err
		}
		checker.acceptedOnL2.Push(acceptedOnL2{
			TransactionHash: tx.Hash,
			Height:          tx.Height,
		})
		return nil
	}
	if block.DeployAccountCount > 0 {
		tx, err := byHeight[storage.DeployAccount, storage.DeployAccountFilter](ctx, checker.deployAccounts, storage.DeployAccountFilter{
			Height: storage.IntegerFilter{
				Eq: block.Height,
			},
		})
		if err != nil {
			return err
		}
		checker.acceptedOnL2.Push(acceptedOnL2{
			TransactionHash: tx.Hash,
			Height:          tx.Height,
		})
		return nil
	}
	if block.DeclareCount > 0 {
		tx, err := byHeight[storage.Declare, storage.DeclareFilter](ctx, checker.declares, storage.DeclareFilter{
			Height: storage.IntegerFilter{
				Eq: block.Height,
			},
		})
		if err != nil {
			return err
		}
		checker.acceptedOnL2.Push(acceptedOnL2{
			TransactionHash: tx.Hash,
			Height:          tx.Height,
		})
		return nil
	}
	if block.L1HandlerCount > 0 {
		tx, err := byHeight[storage.L1Handler, storage.L1HandlerFilter](ctx, checker.l1Handlers, storage.L1HandlerFilter{
			Height: storage.IntegerFilter{
				Eq: block.Height,
			},
		})
		if err != nil {
			return err
		}
		checker.acceptedOnL2.Push(acceptedOnL2{
			TransactionHash: tx.Hash,
			Height:          tx.Height,
		})
		return nil
	}

	checker.acceptedOnL2.Push(acceptedOnL2{
		Height: block.Height,
	})

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
	checker.acceptedOnL2.Push(newAcceptedOnL2FromIndexingBlock(block))
}

func (checker *statusChecker) getStatus(ctx context.Context, item acceptedOnL2) (storage.Status, error) {
	if len(item.TransactionHash) > 0 {
		return checker.receiver.TransactionStatus(ctx, encoding.EncodeHex(item.TransactionHash))
	}

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
