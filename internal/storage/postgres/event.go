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

// Event -
type Event struct {
	*postgres.Table[*storage.Event]
}

// NewEvent -
func NewEvent(db *database.PgGo) *Event {
	return &Event{
		Table: postgres.NewTable[*storage.Event](db),
	}
}

// InsertByCopy -
func (event *Event) InsertByCopy(events []storage.Event) (io.Reader, string, error) {
	if len(events) == 0 {
		return nil, "", nil
	}

	builder := new(strings.Builder)

	for i := range events {
		if err := writeUint64(builder, events[i].ID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, events[i].Height); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeTime(builder, events[i].Time); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].DeclareID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].DeployID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].DeployAccountID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].InvokeID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].L1HandlerID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].FeeID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64Pointer(builder, events[i].InternalID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, events[i].Order); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, events[i].ContractID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeUint64(builder, events[i].FromID); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, events[i].Keys...); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, events[i].Data...); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeString(builder, events[i].Name); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeMap(builder, events[i].ParsedData); err != nil {
			return nil, "", err
		}

		if err := builder.WriteByte('\n'); err != nil {
			return nil, "", err
		}
	}

	query := fmt.Sprintf(`COPY %s FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.Event{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (event *Event) Filter(ctx context.Context, fltr storage.EventFilter, opts ...storage.FilterOption) ([]storage.Event, error) {
	q := event.DB().ModelContext(ctx, (*storage.Event)(nil))
	q = integerFilter(q, "id", fltr.ID)
	q = integerFilter(q, "height", fltr.Height)
	q = timeFilter(q, "time", fltr.Time)
	q = idFilter(q, "contract_id", fltr.Contract, "Contract")
	q = idFilter(q, "from_id", fltr.From, "From")
	q = stringFilter(q, "name", fltr.Name)
	q = jsonFilter(q, "parsed_data", fltr.ParsedData)
	q = optionsFilter(q, opts...)

	var result []storage.Event
	err := q.Select(&result)
	return result, err
}
