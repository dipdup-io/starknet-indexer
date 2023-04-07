package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/dipdup-io/starknet-go-api/pkg/encoding"
	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
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
	wg             *sync.WaitGroup
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
		wg:             new(sync.WaitGroup),
	}
}

// Start -
func (checker *statusChecker) Start(ctx context.Context) {
	checker.wg.Add(1)
	go checker.start(ctx)
}

func (checker *statusChecker) start(ctx context.Context) {
	defer checker.wg.Done()

	if err := checker.init(ctx); err != nil {
		log.Err(err).Msg("checker init")
		return
	}

	if err := checker.check(ctx); err != nil {
		log.Err(err).Msg("check block status")
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := checker.check(ctx); err != nil {
				log.Err(err).Msg("check block status")
			}
		}
	}
}

// Close -
func (checker *statusChecker) Close() error {
	checker.wg.Wait()
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

func byHeight[T sdk.Model](ctx context.Context, src storage.Heightable[T], height uint64) (t T, err error) {
	tx, err := src.ByHeight(ctx, height, 1, 0)
	if err != nil {
		return t, err
	}
	if len(tx) == 0 {
		return t, errors.Wrapf(
			errCantFindTransactionsInBlock,
			"model = %s height = %d", t.TableName(), height)
	}

	return tx[0], nil
}

func (checker *statusChecker) addIndexedBlockToQueue(ctx context.Context, block storage.Block) error {
	if block.InvokeCount > 0 {
		tx, err := byHeight[storage.Invoke](ctx, checker.invoke, block.Height)
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
		tx, err := byHeight[storage.Deploy](ctx, checker.deploys, block.Height)
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
		tx, err := byHeight[storage.DeployAccount](ctx, checker.deployAccounts, block.Height)
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
		tx, err := byHeight[storage.Declare](ctx, checker.declares, block.Height)
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
		tx, err := byHeight[storage.L1Handler](ctx, checker.l1Handlers, block.Height)
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

		log.Info().Str("status", status.String()).Uint64("height", item.Height).Msg("update block status")
	}
}

func (checker *statusChecker) addBlock(block storage.Block) {
	checker.acceptedOnL2.Push(newAcceptedOnL2FromIndexingBlock(block))
}

func (checker *statusChecker) getStatus(ctx context.Context, item acceptedOnL2) (storage.Status, error) {
	if len(item.TransactionHash) > 0 {
		return checker.receiver.TransactionStatus(ctx, encoding.EncodeHex(item.TransactionHash))
	}

	block, err := checker.receiver.GetBlock(ctx, item.Height)
	if err != nil {
		return storage.StatusUnknown, err
	}

	status := storage.NewStatus(block.Status)
	return status, nil
}

func (checker *statusChecker) update(ctx context.Context, height uint64, status storage.Status) error {
	tx, err := checker.transactable.BeginTransaction(ctx)
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
		if _, err := tx.Exec(ctx, `update ? set status = ? where height = ?`, pg.Ident(model.TableName()), status, height); err != nil {
			return tx.HandleError(ctx, err)
		}
	}

	return tx.Flush(ctx)
}
