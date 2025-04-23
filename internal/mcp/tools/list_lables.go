package tools

import (
	"context"
	"encoding/json"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pkg/errors"
)

type Table struct {
	Name string `bun:"table_name" json:"table_name"`
}

func ListTablesTool(ctx context.Context, s postgres.Storage, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := `
	   SELECT DISTINCT
	       regexp_replace(table_name, '_\d{4}_\d{2}$', '') AS table_name
	   FROM information_schema.tables
	   WHERE table_schema = 'public'
	   ORDER BY table_name`

	var tables []Table
	err := s.Connection().DB().NewRaw(query).Scan(ctx, &tables)

	tableNames := make([]string, len(tables))
	for i := range tables {
		tableNames[i] = tables[i].Name
	}

	jsonTableNames, err := json.Marshal(tableNames)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshalling table names")
	}
	return mcp.NewToolResultText(string(jsonTableNames)), nil
}
