package tools

import (
	"context"
	"encoding/json"
	"github.com/dipdup-io/starknet-indexer/internal/mcp/db"
	"github.com/dipdup-net/go-lib/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

func ListTablesTool(dbConfig config.Database, _ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sqlDB, err := db.CreateDBConnection(dbConfig)
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	database := &db.Database{Db: sqlDB}

	query := `
	   SELECT DISTINCT
	       regexp_replace(table_name, '_\d{4}_\d{2}$', '') AS table_name
	   FROM information_schema.tables
	   WHERE table_schema = 'public'
	   ORDER BY table_name`

	results, err := database.ExecuteQuery(query, make(map[string]any))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute list tables")
	}

	var tableNames []string
	for _, row := range results {
		if tableName, ok := row["table_name"].(string); ok {
			tableNames = append(tableNames, tableName)
		}
	}

	jsonTables, err := json.Marshal(tableNames)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling table names")
	}

	return mcp.NewToolResultText(string(jsonTables)), nil
}
