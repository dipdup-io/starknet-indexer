package sequencer

import (
	"context"
	"github.com/dipdup-io/starknet-indexer/pkg/indexer/sqd_receiver/api"
)

func (s *Module) listen(ctx context.Context) {
	s.Log.Info().Msg("module started")

	input := s.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				s.Log.Warn().Msg("can't read message from input, it was drained and closed")
				s.MustOutput(StopOutput).Push(struct{}{})
				return
			}
			block, ok := msg.(*api.SqdBlockResponse)
			if !ok {
				s.Log.Warn().Msgf("invalid message type: %T", msg)
				continue
			}

			s.buffer[block.Header.Number] = block

			s.Log.Info().
				Uint64("ID", block.Header.Number).
				Msg("received block")
		}
	}
}
