package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// L1Handler -
type L1Handler struct {
	*postgres.Table[*storage.L1Handler]
}

// NewL1Handler -
func NewL1Handler(db *database.Bun) *L1Handler {
	return &L1Handler{
		Table: postgres.NewTable[*storage.L1Handler](db),
	}
}

// Filter -
func (l1 *L1Handler) Filter(ctx context.Context, fltr []storage.L1HandlerFilter, opts ...storage.FilterOption) (result []storage.L1Handler, err error) {
	query := l1.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" AND ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "l1_handler.id", fltr[i].ID)
				q = integerFilter(q, "l1_handler.height", fltr[i].Height)
				q = timeFilter(q, "l1_handler.time", fltr[i].Time)
				q = enumFilter(q, "l1_handler.status", fltr[i].Status)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = equalityFilter(q, "l1_handler.selector", fltr[i].Selector)
				q = stringFilter(q, "l1_handler.entrypoint", fltr[i].Entrypoint)
				q = jsonFilter(q, "l1_handler.parsed_calldata", fltr[i].ParsedCalldata)
				return q
			})
		}
		return q1
	})

	query = optionsFilter(query, "l1_handler", opts...)
	query.Relation("Contract")

	err = query.Scan(ctx)
	return result, err
}

func (l1 *L1Handler) HashByHeight(ctx context.Context, height uint64) (hash []byte, err error) {
	err = l1.DB().NewSelect().
		Model((*storage.L1Handler)(nil)).
		Column("hash").
		Where("height = ?", height).
		Limit(1).
		Scan(ctx, &hash)
	return
}
