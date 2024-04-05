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
	query := t.DB().NewSelect().Model((*storage.Transfer)(nil))
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

	err = t.DB().NewSelect().TableExpr("(?) as transfer", query).
		ColumnExpr("transfer.*").
		ColumnExpr("from_addr.id as from__id, from_addr.class_id as from__class_id, from_addr.hash as from__hash, from_addr.height as from__height").
		ColumnExpr("to_addr.id as to__id, to_addr.class_id as to__class_id, to_addr.hash as to__hash, to_addr.height as to__height").
		ColumnExpr("contract.id as contract__id, contract.class_id as contract__class_id, contract.hash as contract__hash, contract.height as contract__height").
		Join("left join address as from_addr on from_addr.id = transfer.from_id").
		Join("left join address as to_addr on to_addr.id = transfer.to_id").
		Join("left join address as contract on contract.id = transfer.contract_id").
		Scan(ctx, &result)
	return result, err
}
