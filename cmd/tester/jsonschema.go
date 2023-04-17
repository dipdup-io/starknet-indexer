package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/dipdup-io/starknet-go-api/pkg/abi"
	models "github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-io/starknet-indexer/internal/storage/postgres"
	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/rs/zerolog/log"
	"github.com/tfkhsr/jsonschema"
)

// JsonSchemaTester -
type JsonSchemaTester struct {
	postgres postgres.Storage
	wg       *sync.WaitGroup
}

// NewJsonSchemaTester -
func NewJsonSchemaTester(postgres postgres.Storage) JsonSchemaTester {
	return JsonSchemaTester{
		postgres: postgres,
		wg:       new(sync.WaitGroup),
	}
}

// String -
func (js JsonSchemaTester) String() string {
	return "json schema tester"
}

// Test -
func (js JsonSchemaTester) Test(ctx context.Context) error {
	log.Info().Msg("start testing json schema...")
	var (
		limit  = 100
		offset = 0
		end    = false
	)
	for !end {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		classes, err := js.postgres.Class.List(ctx, uint64(limit), uint64(offset), storage.SortOrderAsc)
		if err != nil {
			return err
		}

		for i := range classes {
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			js.wg.Add(1)
			go func(c *models.Class) {
				defer js.wg.Done()

				var a abi.Abi
				if err := json.Unmarshal(c.Abi, &a); err != nil {
					log.Err(err).Msgf("can't unmarshal abi %x", c.Hash)
					return
				}

				schema := a.JsonSchema()
				b, err := json.MarshalIndent(schema, "", " ")
				if err != nil {
					log.Err(err).Msgf("can't marshal json schema %x", c.Hash)
					return
				}

				if _, err := jsonschema.Parse(b); err != nil {
					log.Err(err).Msgf("invalid json schema %s", c.Hash)
					return
				}

				log.Info().Msgf("json schema of class %x is valid", c.Hash)
			}(classes[i])
		}

		offset += len(classes)
		end = len(classes) < limit
	}

	js.wg.Wait()

	return nil
}

// Close -
func (js JsonSchemaTester) Close() error {
	js.wg.Wait()

	return nil
}
