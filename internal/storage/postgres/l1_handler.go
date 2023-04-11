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

// L1Handler -
type L1Handler struct {
	*postgres.Table[*storage.L1Handler]
}

// NewL1Handler -
func NewL1Handler(db *database.PgGo) *L1Handler {
	return &L1Handler{
		Table: postgres.NewTable[*storage.L1Handler](db),
	}
}

// InsertByCopy -
func (l1 *L1Handler) InsertByCopy(txs []storage.L1Handler) (io.Reader, string, error) {
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
		if err := writeUint64(builder, uint64(txs[i].Status)); err != nil {
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
		if err := writeUint64(builder, txs[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, uint64(txs[i].Position)); err != nil {
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
		if err := writeDecimal(builder, txs[i].MaxFee); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeDecimal(builder, txs[i].Nonce); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, txs[i].CallData...); err != nil {
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

	query := fmt.Sprintf(`COPY %s FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.L1Handler{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (l1 *L1Handler) Filter(ctx context.Context, fltr storage.L1HandlerFilter, opts ...storage.FilterOption) ([]storage.L1Handler, error) {
	q := l1.DB().ModelContext(ctx, (*storage.L1Handler)(nil))
	q = integerFilter(q, "id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = enumFilter(q, "status", fltr.Status)
	q = addressFilter(q, "hash", fltr.Contract, "Contract")
	q = equalityFilter(q, "selector", fltr.Selector)
	q = stringFilter(q, "entrypoint", fltr.Entrypoint)
	q = jsonFilter(q, "parsed_calldata", fltr.ParsedCalldata)
	q = optionsFilter(q, opts...)

	var result []storage.L1Handler
	err := q.Select(&result)
	return result, err
}
