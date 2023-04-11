package grpc

import (
	"context"
	"math"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	generalPB "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/pb"
	"golang.org/x/exp/slices"
)

const (
	priorityTokenBalance  = 21
	priorityStorageDiff   = 22
	priorityTransfer      = 23
	priorityMessage       = 24
	priorityEvent         = 25
	priorityInternal      = 26
	priorityFee           = 27
	priorityDeployAccount = 28
	priorityDeploy        = 29
	priorityDeclare       = 30
	priorityL1Handler     = 31
	priorityInvoke        = 32
)

type table[T storage.Heightable, F any] struct {
	Data         []T
	store        storage.Filterable[T, F]
	fltr         F
	priority     int
	offset       int
	limit        int
	targetHeight uint64

	end bool
}

func newTable[T storage.Heightable, F any](store storage.Filterable[T, F], fltr F, priority int) *table[T, F] {
	return &table[T, F]{
		Data:     make([]T, 0),
		store:    store,
		fltr:     fltr,
		priority: priority,
		limit:    100,
	}
}

func (t *table[T, F]) getFirst() (T, bool) {
	if len(t.Data) == 0 {
		var m T
		return m, false
	}

	return t.Data[0], true
}

// Push -
func (t *table[T, F]) Push(arr []T) {
	t.Data = append(t.Data, arr...)
}

// Pop -
func (t *table[T, F]) Pop() (T, bool) {
	result, ok := t.getFirst()
	if ok {
		t.Data = t.Data[1:]
		t.end = len(t.Data) == 0 || t.targetHeight < t.GetHeight()
	}
	return result, ok
}

// GetHeight -
func (t *table[T, F]) GetHeight() uint64 {
	if result, ok := t.getFirst(); ok {
		return result.GetHeight()
	}
	return math.MaxUint64
}

// Priority -
func (t *table[T, F]) Priority() int {
	return t.priority
}

// Len -
func (t *table[T, F]) Len() int {
	return len(t.Data)
}

// IsFinished -
func (t *table[T, F]) IsFinished() bool {
	return t.end
}

// Head -
func (t *table[T, F]) Head() any {
	if head, ok := t.Pop(); ok {
		return head
	}
	return nil
}

// SetTargetHeight -
func (t *table[T, F]) SetTargetHeight(targetHeight uint64) {
	t.targetHeight = targetHeight
}

// Receive -
func (t *table[T, F]) Receive(ctx context.Context) error {
	if t.end || t.Len() > 0 {
		return nil
	}

	data, err := t.store.Filter(
		ctx,
		t.fltr,
		storage.WithAscSortByIdFilter(),
		storage.WithLimitFilter(t.limit),
		storage.WithOffsetFilter(t.offset),
	)
	if err != nil {
		return err
	}
	t.Push(data)
	t.offset += len(data)
	t.end = len(data) == 0 || t.GetHeight() > t.targetHeight
	return nil
}

// Reset -
func (t *table[T, F]) Reset() {
	t.end = false
}

type synchronizable interface {
	GetHeight() uint64
	Priority() int
	Receive(ctx context.Context) error
	Len() int
	Head() any
	IsFinished() bool
	SetTargetHeight(targetHeight uint64)
	Reset()
}

type tables []synchronizable

// Finished -
func (a tables) Finished() bool {
	for i := range a {
		if !a[i].IsFinished() {
			return false
		}
	}
	return true
}

func (module *Server) sync(ctx context.Context, subscriptionID uint64, req *pb.SubscribeRequest, stream pb.IndexerService_SubscribeServer) error {
	sf := newSubscriptionFilters(req)

	syncTables := make(tables, 0)
	if sf.invoke != nil {
		syncTables = append(syncTables, newTable[storage.Invoke, storage.InvokeFilter](module.db.Invoke, *sf.invoke, priorityInvoke))
	}
	if sf.l1Handler != nil {
		syncTables = append(syncTables, newTable[storage.L1Handler, storage.L1HandlerFilter](module.db.L1Handler, *sf.l1Handler, priorityL1Handler))
	}
	if sf.declare != nil {
		syncTables = append(syncTables, newTable[storage.Declare, storage.DeclareFilter](module.db.Declare, *sf.declare, priorityDeclare))
	}
	if sf.deploy != nil {
		syncTables = append(syncTables, newTable[storage.Deploy, storage.DeployFilter](module.db.Deploy, *sf.deploy, priorityDeploy))
	}
	if sf.deployAccount != nil {
		syncTables = append(syncTables, newTable[storage.DeployAccount, storage.DeployAccountFilter](module.db.DeployAccount, *sf.deployAccount, priorityDeployAccount))
	}
	if sf.internal != nil {
		syncTables = append(syncTables, newTable[storage.Internal, storage.InternalFilter](module.db.Internal, *sf.internal, priorityInternal))
	}
	if sf.fee != nil {
		syncTables = append(syncTables, newTable[storage.Fee, storage.FeeFilter](module.db.Fee, *sf.fee, priorityFee))
	}
	if sf.event != nil {
		syncTables = append(syncTables, newTable[storage.Event, storage.EventFilter](module.db.Event, *sf.event, priorityEvent))
	}
	if sf.message != nil {
		syncTables = append(syncTables, newTable[storage.Message, storage.MessageFilter](module.db.Message, *sf.message, priorityMessage))
	}
	if sf.transfer != nil {
		syncTables = append(syncTables, newTable[storage.Transfer, storage.TransferFilter](module.db.Transfer, *sf.transfer, priorityTransfer))
	}
	if sf.storageDiff != nil {
		syncTables = append(syncTables, newTable[storage.StorageDiff, storage.StorageDiffFilter](module.db.StorageDiff, *sf.storageDiff, priorityStorageDiff))
	}

	var height uint64

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		block, err := module.db.Blocks.Last(ctx)
		if err != nil {
			return err
		}
		if height == block.Height {
			break
		}
		height = block.Height

		if err := module.syncTables(ctx, syncTables, height, subscriptionID, stream); err != nil {
			return err
		}
	}

	if sf.tokenBalance != nil {
		if err := module.syncTokenBalances(ctx, *sf.tokenBalance, subscriptionID, stream); err != nil {
			return err
		}
	}

	return nil
}

func (module *Server) syncTables(ctx context.Context, tables tables, targetHeight, subscriptionID uint64, stream pb.IndexerService_SubscribeServer) error {
	for i := range tables {
		tables[i].Reset()
		if err := tables[i].Receive(ctx); err != nil {
			return err
		}
		tables[i].SetTargetHeight(targetHeight)
	}

	for !tables.Finished() {
		slices.SortFunc(tables, func(a synchronizable, b synchronizable) bool {
			aH := a.GetHeight()
			bH := b.GetHeight()
			if aH == bH {
				return a.Priority() > b.Priority()
			}
			return aH < bH
		})
		if head := tables[0].Head(); head != nil {
			if err := sendModelToClient(ctx, subscriptionID, stream, head); err != nil {
				return err
			}
		}

		if tables[0].Len() == 0 && !tables[0].IsFinished() {
			if err := tables[0].Receive(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func sendModelToClient(ctx context.Context, subscriptionID uint64, stream pb.IndexerService_SubscribeServer, model any) error {
	var msg pb.Subscription
	switch typ := model.(type) {
	case storage.Invoke:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Invoke: Invoke(&typ),
		}
	case storage.L1Handler:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			L1Handler: L1Handler(&typ),
		}
	case storage.Declare:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Declare: Declare(&typ),
		}
	case storage.Deploy:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Deploy: Deploy(&typ),
		}
	case storage.DeployAccount:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			DeployAccount: DeployAccount(&typ),
		}
	case storage.Internal:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Internal: Internal(&typ),
		}
	case storage.Fee:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Fee: Fee(&typ),
		}
	case storage.Event:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Event: Event(&typ),
		}
	case storage.Message:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Message: Message(&typ),
		}
	case storage.StorageDiff:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			StorageDiff: StorageDiff(&typ),
		}
	case storage.Transfer:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			Transfer: Transfer(&typ),
		}
	case storage.TokenBalance:
		msg = pb.Subscription{
			Response: &generalPB.SubscribeResponse{
				Id: subscriptionID,
			},
			TokenBalance: TokenBalance(&typ),
		}
	default:
		return nil
	}
	return stream.Send(&msg)
}

func (module *Server) syncTokenBalances(ctx context.Context, fltr storage.TokenBalanceFilter, subscriptionID uint64, stream pb.IndexerService_SubscribeServer) error {
	var (
		offset int
		end    bool
		limit  = 100
	)

	for !end {
		data, err := module.db.TokenBalance.Filter(ctx, fltr, storage.WithLimitFilter(limit), storage.WithOffsetFilter(offset))
		if err != nil {
			return err
		}
		end = len(data) < limit
		offset += len(data)
		for i := range data {
			if err := sendModelToClient(ctx, subscriptionID, stream, data[i]); err != nil {
				return err
			}
		}
	}
	return nil
}