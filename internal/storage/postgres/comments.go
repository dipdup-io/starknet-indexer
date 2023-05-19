package postgres

import (
	"context"
	"reflect"
	"strings"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/go-lib/hasura"
	sdk "github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

func makeComments(ctx context.Context, conn *database.PgGo) error {
	for _, model := range []sdk.Model{
		storage.Address{},
		storage.Block{},
		storage.Class{},
		storage.Declare{},
		storage.DeployAccount{},
		storage.Deploy{},
		storage.Token{},
		storage.Event{},
		storage.Fee{},
		storage.Internal{},
		storage.Invoke{},
		storage.L1Handler{},
		storage.Message{},
		storage.Proxy{},
		storage.TokenBalance{},
		storage.Transfer{},
		storage.State{},
		storage.StorageDiff{},
	} {
		typ := reflect.TypeOf(model)
		for i := 0; i < typ.NumField(); i++ {
			fieldType := typ.Field(i)
			pgTag, ok := fieldType.Tag.Lookup("pg")
			if !ok {
				continue
			}

			tags := strings.Split(pgTag, ",")

			var name string
			for i := range tags {
				if i == 0 {
					if name == "" {
						name = hasura.ToSnakeCase(fieldType.Name)
					} else {
						name = tags[i]
					}
					continue
				}

				parts := strings.Split(tags[i], ":")
				if parts[0] == "comment" {
					if len(parts) != 2 {
						return errors.Errorf("invalid comments format: %s", pgTag)
					}
					if fieldType.Name == "tableName" {
						if _, err := conn.DB().ExecContext(ctx, `COMMENT ON TABLE ? IS ?`, pg.Safe(model.TableName()), parts[1]); err != nil {
							return err
						}
					} else {
						if _, err := conn.DB().ExecContext(ctx, `COMMENT ON COLUMN ?.? IS ?`, pg.Safe(model.TableName()), pg.Safe(name), parts[1]); err != nil {
							return err
						}
					}
					continue
				}
			}
		}
	}

	return nil
}
