package postgres

import (
	"context"
	"fmt"
	"time"

	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
)

// PartitionManager -
type PartitionManager struct {
	conn *database.PgGo

	lastId string
}

// NewPartitionManager -
func NewPartitionManager(conn *database.PgGo) PartitionManager {
	return PartitionManager{
		conn: conn,
	}
}

const createPartitionTemplate = `CREATE TABLE IF NOT EXISTS ? PARTITION OF ? FOR VALUES FROM (?) TO (?);`

func boundaries(current time.Time) (time.Time, time.Time) {
	start := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return start, end
}

func (pm *PartitionManager) partitionId(currentTime time.Time) string {
	return fmt.Sprintf("%d_%02d", currentTime.Year(), currentTime.Month())
}

// CreatePartitions -
func (pm *PartitionManager) CreatePartitions(ctx context.Context, currentTime time.Time) error {
	id := pm.partitionId(currentTime)
	if id == pm.lastId {
		return nil
	}

	start, end := boundaries(currentTime)

	for _, model := range []storage.Model{
		&models.Internal{},
		&models.Declare{},
		&models.DeployAccount{},
		&models.Deploy{},
		&models.Event{},
		&models.Invoke{},
		&models.L1Handler{},
		&models.Message{},
		&models.Transfer{},
		&models.Fee{},
	} {
		partitionName := fmt.Sprintf("%s_%s", model.TableName(), id)
		if _, err := pm.conn.DB().ExecContext(
			ctx,
			createPartitionTemplate,
			pg.Ident(partitionName),
			pg.Ident(model.TableName()),
			start.Format(time.RFC3339Nano),
			end.Format(time.RFC3339Nano),
		); err != nil {
			return err
		}
	}

	pm.lastId = id
	return nil
}
