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

// Internal -
type Internal struct {
	*postgres.Table[*storage.Internal]
}

// NewInternal -
func NewInternal(db *database.PgGo) *Internal {
	return &Internal{
		Table: postgres.NewTable[*storage.Internal](db),
	}
}

// InsertByCopy -
func (db Internal) InsertByCopy(txs []storage.Internal) (io.Reader, string, error) {
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
		if err := writeUint64Pointer(builder, txs[i].InvokeID); err != nil {
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
		if err := writeUint64Pointer(builder, txs[i].InternalID); err != nil {
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
		if err := writeUint64(builder, txs[i].CallerID); err != nil {
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
		if err := writeUint64(builder, uint64(txs[i].Status)); err != nil {
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
		if err := writeUint64(builder, uint64(txs[i].EntrypointType)); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeBytes(builder, txs[i].Hash); err != nil {
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
		if err := writeStringArray(builder, txs[i].Result...); err != nil {
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
		if err := writeMap(builder, txs[i].ParsedCalldata); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeMap(builder, txs[i].ParsedResult); err != nil {
			return nil, "", err
		}

		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.Internal{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (d *Internal) Filter(ctx context.Context, fltr []storage.InternalFilter, opts ...storage.FilterOption) ([]storage.Internal, error) {
	query := d.DB().ModelContext(ctx, (*storage.Internal)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "internal_tx.id", fltr[i].ID)
				q = integerFilter(q, "internal_tx.height", fltr[i].Height)
				q = timeFilter(q, "internal_tx.time", fltr[i].Time)
				q = enumFilter(q, "internal_tx.status", fltr[i].Status)
				q = addressFilter(q, "internal_tx.contract_id", fltr[i].Contract, "Contract")
				q = addressFilter(q, "internal_tx.caller_id", fltr[i].Caller, "Caller")
				q = addressFilter(q, "internal_tx..class_id", fltr[i].Class, "Class")
				q = equalityFilter(q, "internal_tx.selector", fltr[i].Selector)
				q = stringFilter(q, "internal_tx.entrypoint", fltr[i].Entrypoint)
				q = enumFilter(q, "internal_tx.entrypoint_type", fltr[i].EntrypointType)
				q = enumFilter(q, "internal_tx.call_type", fltr[i].CallType)
				q = jsonFilter(q, "internal_tx.parsed_calldata", fltr[i].ParsedCalldata)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "internal_tx", opts...)
	query.Relation("Contract").Relation("Caller").Relation("Class")

	var result []storage.Internal
	err := query.Select(&result)
	return result, err
}
