package postgres

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
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
func (t *Transfer) Filter(ctx context.Context, fltr storage.TransferFilter, opts ...storage.FilterOption) ([]storage.Transfer, error) {
	q := t.DB().ModelContext(ctx, (*storage.Transfer)(nil))
	q = integerFilter(q, "transfer.id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = addressFilter(q, "hash", fltr.Contract, "Contract")
	q = addressFilter(q, "hash", fltr.From, "From")
	q = addressFilter(q, "hash", fltr.To, "To")
	q = stringFilter(q, "token_id", fltr.TokenId)
	q = optionsFilter(q, opts...)

	var result []storage.Transfer
	err := q.Select(&result)
	return result, err
}
