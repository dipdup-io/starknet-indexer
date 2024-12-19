package adapter

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Adapter struct {
	modules.BaseModule
}

var _ modules.Module = (*Adapter)(nil)

const (
	InputName  = "blocks"
	OutputName = "parsed_blocks"
	StopOutput = "stop"
)

func New() Adapter {
	m := Adapter{
		BaseModule: modules.New("sqd adapter"),
	}
	m.CreateInputWithCapacity(InputName, 128)
	m.CreateOutput(OutputName)
	m.CreateOutput(StopOutput)

	return m
}

func (a *Adapter) Start(ctx context.Context) {
	a.Log.Info().Msg("starting...")
	a.G.GoCtx(ctx, a.listen)
}

func (a *Adapter) Close() error {
	a.Log.Info().Msg("closing...")
	a.G.Wait()
	return nil
}
