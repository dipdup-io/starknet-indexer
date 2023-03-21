package postgres

import (
	"context"

	"github.com/dipdup-io/starknet-indexer/internal/storage"
	"github.com/dipdup-net/go-lib/database"
	"github.com/dipdup-net/indexer-sdk/pkg/storage/postgres"
)

// State -
type State struct {
	*postgres.Table[*storage.State]
}

// NewState -
func NewState(db *database.PgGo) *State {
	return &State{
		Table: postgres.NewTable[*storage.State](db),
	}
}

// ByName -
func (s *State) ByName(ctx context.Context, name string) (state storage.State, err error) {
	err = s.DB().ModelContext(ctx, &state).Where("name = ?", name).Select(&state)
	return
}
