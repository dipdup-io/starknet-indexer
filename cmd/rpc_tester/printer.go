package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dipdup-io/starknet-indexer/pkg/grpc/pb"
	"github.com/dipdup-net/indexer-sdk/pkg/modules/printer"
	"github.com/rs/zerolog/log"
)

// Printer -
type Printer struct {
	*printer.Printer

	eventCounters map[string]*atomic.Uint64

	wg *sync.WaitGroup
}

// NewPrinter -
func NewPrinter() *Printer {
	return &Printer{
		printer.NewPrinter(),
		make(map[string]*atomic.Uint64),
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
					p.incrementEventCounter(typ.Event.Name)

					log.Info().
						Str("name", typ.Event.Name).
						Uint64("height", typ.Event.Height).
						Uint64("time", typ.Event.Time).
						Uint64("id", typ.Event.Id).
						Uint64("subscription", typ.Response.Id).
						Str("contract", fmt.Sprintf("0x%x", typ.Event.Contract)).
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

	for name, counter := range p.eventCounters {
		log.Info().Str("event", name).Uint64("count", counter.Load()).Msg("events were handled")
	}

	if err := p.Printer.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Printer) incrementEventCounter(name string) {
	if counter, ok := p.eventCounters[name]; ok {
		counter.Add(1)
	} else {
		counter := new(atomic.Uint64)
		counter.Add(1)
		p.eventCounters[name] = counter
	}
}
