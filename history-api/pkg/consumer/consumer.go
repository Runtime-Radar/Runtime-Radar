package consumer

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/database/clickhouse"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model/convert"
	"github.com/runtime-radar/runtime-radar/lib/rabbit"
)

type Consumer struct {
	PublishConsumer        rabbit.PublishConsumer
	RuntimeEventRepository clickhouse.RuntimeEventRepository
}

func (c *Consumer) Run(stop <-chan struct{}) {
	log.Info().Msgf("Runtime events consumer started")
	defer log.Info().Msgf("Runtime events consumer stopped")

	for {
		select {
		default:
			ev := &api.RuntimeEvent{}
			if err := c.PublishConsumer.Consume(context.Background(), ev); err != nil {
				log.Error().Msgf("Can't consume runtime event: %v", err)
				continue
			}

			m, err := convert.RuntimeEventFromProto(ev)
			if err != nil {
				log.Error().Err(err).Msg("Can't convert runtime event to model")
				continue
			}

			if err := c.RuntimeEventRepository.Add(context.Background(), &[]model.RuntimeEvent{m}); err != nil {
				log.Error().Err(err).Msgf("Can't save runtime event")
			}

		case <-stop:
			return
		}
	}
}
