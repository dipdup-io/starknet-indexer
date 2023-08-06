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

// Deploy -
type Deploy struct {
	*postgres.Table[*storage.Deploy]
}

// NewDeploy -
func NewDeploy(db *database.PgGo) *Deploy {
	return &Deploy{
		Table: postgres.NewTable[*storage.Deploy](db),
	}
}

// InsertByCopy -
func (d *Deploy) InsertByCopy(txs []storage.Deploy) (io.Reader, string, error) {
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
		if err := writeUint64(builder, txs[i].ClassID); err != nil {
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
		if err := writeBytes(builder, txs[i].ContractAddressSalt); err != nil {
			return nil, "", err
		}
		if err := builder.WriteByte(','); err != nil {
			return nil, "", err
		}
		if err := writeStringArray(builder, txs[i].ConstructorCalldata...); err != nil {
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

	query := fmt.Sprintf(`COPY %s FROM STDIN WITH (FORMAT csv, ESCAPE '\', QUOTE '"', DELIMITER ',')`, storage.Deploy{}.TableName())
	return strings.NewReader(builder.String()), query, nil
}

// Filter -
func (d *Deploy) Filter(ctx context.Context, fltr []storage.DeployFilter, opts ...storage.FilterOption) ([]storage.Deploy, error) {
	query := d.DB().ModelContext(ctx, (*storage.Deploy)(nil))
	query = query.WhereGroup(func(q1 *orm.Query) (*orm.Query, error) {
		for i := range fltr {
			q1 = q1.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = integerFilter(q, "deploy.id", fltr[i].ID)
				q = integerFilter(q, "deploy.height", fltr[i].Height)
				q = timeFilter(q, "deploy.time", fltr[i].Time)
				q = enumFilter(q, "deploy.status", fltr[i].Status)
				q = addressFilter(q, "deploy.class_id", fltr[i].Class, "Class")
				q = jsonFilter(q, "deploy.parsed_calldata", fltr[i].ParsedCalldata)
				return q, nil
			})
		}
		return q1, nil
	})
	query = optionsFilter(query, "deploy", opts...)
	query.Relation("Contract").Relation("Class")

	var result []storage.Deploy
	err := query.Select(&result)
	return result, err
}
