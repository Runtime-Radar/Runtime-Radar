//go:build !tinygo.wasm

package consumer

import (
	"context"

	"github.com/cilium/tetragon/api/v1/tetragon"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/processor"
	"github.com/runtime-radar/runtime-radar/lib/rabbit"
)

type Consumer struct {
	PublishConsumer rabbit.PublishConsumer
	Processor       processor.Processor
}

func (c *Consumer) Run(stop <-chan struct{}) {
	log.Info().Msgf("Events consumer started")
	defer log.Info().Msgf("Events consumer stopped")

	jobs := c.Processor.Jobs()
	ctx := context.Background()

	for {
		select {
		default:
			ev := &tetragon.GetEventsResponse{}
			if err := c.PublishConsumer.Consume(ctx, ev); err != nil {
				log.Error().Msgf("Can't consume event: %v", err)

				continue
			}

			log.Debug().Interface("event", ev).Msgf("Got event")

			jobs <- ev
		case <-stop:
			return
		}
	}
}
