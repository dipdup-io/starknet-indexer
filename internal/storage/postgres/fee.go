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

// Fee -
type Fee struct {
	*postgres.Table[*storage.Fee]
}

// NewFee -
func NewFee(db *database.PgGo) *Fee {
	return &Fee{
		Table: postgres.NewTable[*storage.Fee](db),
	}
}

// InsertByCopy -
func (fee *Fee) InsertByCopy(txs []storage.Fee) (io.Reader, string, error) {
	if len(txs) == 0 {
		return nil, "", nil
	}
	builder := new(strings.Builder)

	for i := range txs {
		if err := writeUint64(builder, txs[i].ID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, txs[i].Height); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeTime(builder, txs[i].Time); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, txs[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, txs[i].CallerID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, txs[i].ClassID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, txs[i].InvokeID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, txs[i].DeclareID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, txs[i].DeployID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, txs[i].DeployAccountID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, txs[i].L1HandlerID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, uint64(txs[i].EntrypointType)); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, uint64(txs[i].CallType)); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, uint64(txs[i].Status)); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, txs[i].Selector); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeString(builder, txs[i].Entrypoint); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, txs[i].Calldata...); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, txs[i].Result...); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeMap(builder, txs[i].ParsedCalldata); err != nil {
			return nil, "", err
		}

		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.Fee{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (fee *Fee) Filter(ctx context.Context, fltr []storage.FeeFilter, opts ...storage.FilterOption) ([]storage.Fee, error) {
	query := fee.DB().ModelContext(ctx, (*storage.Fee)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "fee.id", fltr[i].ID)
				q = integerFilter(q, "fee.height", fltr[i].Height)
				q = timeFilter(q, "fee.time", fltr[i].Time)
				q = enumFilter(q, "fee.status", fltr[i].Status)
				q = addressFilter(q, "fee.contract_id", fltr[i].Contract, "Contract")
				q = addressFilter(q, "fee.caller_id", fltr[i].Caller, "Caller")
				q = addressFilter(q, "fee.class_id", fltr[i].Class, "Class")
				q = equalityFilter(q, "fee.selector", fltr[i].Selector)
				q = stringFilter(q, "fee.entrypoint", fltr[i].Entrypoint)
				q = enumFilter(q, "fee.entrypoint_type", fltr[i].EntrypointType)
				q = enumFilter(q, "fee.call_type", fltr[i].CallType)
				q = jsonFilter(q, "fee.parsed_calldata", fltr[i].ParsedCalldata)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "fee", opts...)

	var result []storage.Fee
	err := query.Select(&result)
	return result, err
}
