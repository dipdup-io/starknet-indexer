package db

import (
	"database/sql"
	"fmt"
	"github.com/dipdup-net/go-lib/config"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type Database struct {
	Db *sql.DB
}

func CreateDBConnection(dbConfig config.Database) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error during connectiong MCP to DB")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "error pinging DB from MCP")
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(1 * time.Minute)

	return db, nil
}

func (d *Database) ExecuteQuery(query string, params map[string]any) ([]map[string]any, error) {
	stmt, err := d.Db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "error preparing query")
	}
	defer stmt.Close()

	args := make([]any, 0, len(params))
	for _, v := range params {
		args = append(args, v)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing query")
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting columns")
	}

	var results []map[string]any

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, errors.Wrapf(err, "error scanning row")
		}

		entry := make(map[string]any)
		for i, col := range columns {
			entry[col] = values[i]
		}

		results = append(results, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "error after row iteration")
	}

	return results, nil
}
