package sequencer

import (
	"context"
)

func (s *Module) listen(ctx context.Context) {
	s.Log.Info().Msg("module started")

	input := s.MustInput(InputName)

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-input.Listen():
			if !ok {
				s.Log.Warn().Msg("can't read message from input, it was drained and closed")
				s.MustOutput(StopOutput).Push(struct{}{})
				return
			}
		}
	}
}
