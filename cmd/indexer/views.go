package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
)

func createViews(ctx context.Context, strg postgres.Storage) ([]string, error) {
	files, err := os.ReadDir("views")
	if err != nil {
		return nil, err
	}

	views := make([]string, 0)
	for i := range files {
		if files[i].IsDir() {
			continue
		}

		path := fmt.Sprintf("views/%s", files[i].Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if _, err := strg.Connection().DB().Exec(string(raw)); err != nil {
			return nil, err
		}
		views = append(views, strings.Split(files[i].Name(), ".")[0])
	}

	return views, nil
}
