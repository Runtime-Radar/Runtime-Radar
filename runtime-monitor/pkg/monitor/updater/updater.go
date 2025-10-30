package updater

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/runtime-monitor/pkg/database"
	"github.com/runtime-radar/runtime-radar/runtime-monitor/pkg/monitor"
	"github.com/runtime-radar/runtime-radar/runtime-monitor/pkg/monitor/config"
)

type Updater struct {
	Interval         time.Duration
	ConfigRepository database.ConfigRepository
	Monitor          monitor.Monitor
}

func (u *Updater) Run(stop <-chan struct{}) {
	log.Debug().Msgf("Config updater started")
	defer log.Debug().Msgf("Config updater stopped")

	t := time.NewTicker(u.Interval)

	for {
		select {
		case <-t.C:
			ctx := context.Background()

			cfg, err := u.ConfigRepository.GetLast(ctx, true) // preload is on
			if err != nil {
				log.Error().Err(err).Msgf("Can't get last config from DB")
				continue
			}

			oldCfg := u.Monitor.Config()

			log.Debug().Interface("old_config", oldCfg).Msgf("Old monitor config")
			log.Debug().Interface("new_config", cfg).Msgf("New monitor config")

			sel, changed := config.Diff(oldCfg, cfg)
			if changed {
				log.Info().
					Interface("config", cfg).
					Interface("selector", sel).
					Msgf("Monitor config changed, re-initializing")

				u.Monitor.Reinit(sel, cfg)
			} else {
				log.Debug().Msgf("Monitor config didn't change")
			}
		case <-stop:
			return
		}
	}
}
