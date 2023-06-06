package postgres

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
	"github.com/go-pg/pg/v10/orm"
)

// Transfer -
type Transfer struct {
	*postgres.Table[*storage.Transfer]
}

// NewTransfer -
func NewTransfer(db *database.PgGo) *Transfer {
	return &Transfer{
		Table: postgres.NewTable[*storage.Transfer](db),
	}
}

// InsertByCopy -
func (t *Transfer) InsertByCopy(transfers []storage.Transfer) (io.Reader, string, error) {
	if len(transfers) == 0 {
		return nil, "", nil
	}
	builder := new(strings.Builder)

	for i := range transfers {
		if err := writeUint64(builder, transfers[i].Height); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeTime(builder, transfers[i].Time); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, transfers[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, transfers[i].FromID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, transfers[i].ToID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeDecimal(builder, transfers[i].Amount); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeDecimal(builder, transfers[i].TokenID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].InvokeID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].DeclareID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].DeployID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].DeployAccountID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].L1HandlerID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].FeeID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, transfers[i].InternalID); err != nil {
			return nil, "", err
		}

		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s (
		height, time, contract_id, from_id, to_id, amount, token_id, invoke_id, declare_id, deploy_id, deploy_account_id, l1_handler_id, fee_id, internal_id
	) FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.Transfer{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (t *Transfer) Filter(ctx context.Context, fltr []storage.TransferFilter, opts ...storage.FilterOption) ([]storage.Transfer, error) {
	query := t.DB().ModelContext(ctx, (*storage.Transfer)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "transfer.id", fltr[i].ID)
				q = integerFilter(q, "transfer.height", fltr[i].Height)
				q = timeFilter(q, "transfer.time", fltr[i].Time)
				q = addressFilter(q, "transfer.contract_id", fltr[i].Contract, "Contract")
				q = addressFilter(q, "transfer.from_id", fltr[i].From, "From")
				q = addressFilter(q, "transfer.to_id", fltr[i].To, "To")
				q = stringFilter(q, "transfer.token_id", fltr[i].TokenId)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "transfer", opts...)

	var result []storage.Transfer
	err := query.Select(&result)
	return result, err
}
