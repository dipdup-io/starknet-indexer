package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/uptrace/bun"
)

// Transfer -
type Transfer struct {
	*postgres.Table[*storage.Transfer]
}

// NewTransfer -
func NewTransfer(db *database.Bun) *Transfer {
	return &Transfer{
		Table: postgres.NewTable[*storage.Transfer](db),
	}
}

// Filter -
func (t *Transfer) Filter(ctx context.Context, fltr []storage.TransferFilter, opts ...storage.FilterOption) (result []storage.Transfer, err error) {
	query := t.DB().NewSelect().Model(&result)
	query = query.WhereGroup(" OR ", func(q1 *bun.SelectQuery) *bun.SelectQuery {
		for i := range fltr {
			q1 = q1.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				q = integerFilter(q, "transfer.id", fltr[i].ID)
				q = integerFilter(q, "transfer.height", fltr[i].Height)
				q = timeFilter(q, "transfer.time", fltr[i].Time)
				q = addressFilter(q, "hash", fltr[i].Contract, "Contract")
				q = addressFilter(q, "hash", fltr[i].From, "From")
				q = addressFilter(q, "hash", fltr[i].To, "To")
				q = stringFilter(q, "transfer.token_id", fltr[i].TokenId)
				return q
			})
		}
		return q1
	})
	query = optionsFilter(query, "transfer", opts...)
	query.Relation("Contract").Relation("From").Relation("To")

	err = query.Scan(ctx)
	return result, err
}
