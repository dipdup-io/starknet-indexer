package main

import (
	"context"
	"sync"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/printer"
	"github.com/rs/zerolog/log"
)

// Printer -
type Printer struct {
	*printer.Printer

	wg *sync.WaitGroup
}

// NewPrinter -
func NewPrinter() *Printer {
	return &Printer{
		printer.NewPrinter(),
		new(sync.WaitGroup),
	}
}

// Start -
func (p *Printer) Start(ctx context.Context) {
	p.wg.Add(1)
	go p.listen(ctx)
}

func (p *Printer) listen(ctx context.Context) {
	defer p.wg.Done()

	input, err := p.Input(printer.InputName)
	if err != nil {
		log.Err(err).Msg("unknown input")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-input.Listen():
			if !ok {
				continue
			}

			switch typ := msg.(type) {
			case *pb.Subscription:
				switch {
				case typ.GetEndOfBlock():
					log.Info().
						Uint64("subscription", typ.Response.Id).
						Msg("end of block")
				case typ.Event != nil:
					log.Info().
						Str("name", typ.Event.Name).
						Uint64("height", typ.Event.Height).
						Uint64("time", typ.Event.Time).
						Uint64("id", typ.Event.Id).
						Uint64("subscription", typ.Response.Id).
						Msg("event")
				case typ.Block != nil:
					log.Info().
						Uint64("height", typ.Block.Height).
						Msg("new block")
				}
			default:
				log.Info().Msgf("unknown message: %T", typ)
			}
		}
	}
}

// Close -
func (p *Printer) Close() error {
	p.wg.Wait()
	if err := p.Printer.Close(); err != nil {
		return err
	}
	return nil
}
