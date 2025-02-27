package adapter

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/receiver"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Adapter struct {
	modules.BaseModule
	results chan receiver.Result
	head    uint64
}

var _ modules.Module = (*Adapter)(nil)

const (
	BlocksInput  = "blocks"
	HeadInput    = "head"
	BlocksOutput = "parsed_blocks"
	StopOutput   = "stop"
	HeadAchieved = "head_achieved"
)

func New(resultsChannel chan receiver.Result) *Adapter {
	m := &Adapter{
		BaseModule: modules.New("sqd adapter"),
		results:    resultsChannel,
	}
	m.CreateInputWithCapacity(BlocksInput, 128)
	m.CreateInputWithCapacity(HeadInput, 1)
	m.CreateOutput(BlocksOutput)
	m.CreateOutput(StopOutput)
	m.CreateOutput(HeadAchieved)

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
