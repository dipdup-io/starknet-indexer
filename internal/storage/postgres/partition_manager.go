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

	lastYear int
}

// NewPartitionManager -
func NewPartitionManager(conn *database.PgGo) PartitionManager {
	return PartitionManager{
		conn: conn,
	}
}

const createPartitionTemplate = `CREATE TABLE IF NOT EXISTS ? PARTITION OF ? FOR VALUES FROM (?) TO (?);`

// CreatePartitions -
func (pm *PartitionManager) CreatePartitions(ctx context.Context, currentTime time.Time) error {
	year := currentTime.Year()
	if year == pm.lastYear {
		return nil
	}

	firstOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextYear := firstOfYear.AddDate(1, 0, 0)

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
		partitionName := fmt.Sprintf("%s_%d", model.TableName(), year)
		if _, err := pm.conn.DB().ExecContext(
			ctx,
			createPartitionTemplate,
			pg.Ident(partitionName),
			pg.Ident(model.TableName()),
			firstOfYear.Format(time.RFC3339Nano),
			firstOfNextYear.Format(time.RFC3339Nano),
		); err != nil {
			return err
		}
	}

	pm.lastYear = year
	return nil
}
