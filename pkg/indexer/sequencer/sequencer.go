package sequencer

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Module struct {
	modules.BaseModule
}

var _ modules.Module = (*Module)(nil)

const (
	InputName  = "blocks_input"
	OutputName = "blocks_output"
	StopOutput = "stop"
)

func New() Module {
	m := Module{
		BaseModule: modules.New("sequencer"),
	}
	m.CreateInputWithCapacity(InputName, 128)
	m.CreateOutput(OutputName)
	m.CreateOutput(StopOutput)

	return m
}

func (s *Module) Start(ctx context.Context) {
	s.Log.Info().Msg("starting...")
	s.G.GoCtx(ctx, s.listen)
}

func (s *Module) Close() error {
	s.Log.Info().Msg("closing...")
	s.G.Wait()
	return nil
}
