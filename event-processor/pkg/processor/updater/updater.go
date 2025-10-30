//go:build !tinygo.wasm

package updater

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/database"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/processor"
	"github.com/runtime-radar/runtime-radar/event-processor/pkg/processor/config"
)

type Updater struct {
	Interval           time.Duration
	Processor          processor.Processor
	ConfigRepository   database.ConfigRepository
	DetectorRepository database.DetectorRepository
}

func (u *Updater) Run(stop <-chan struct{}) {
	log.Debug().Msgf("Config updater started")
	defer log.Debug().Msgf("Config updater stopped")

	t := time.NewTicker(u.Interval)

	for {
		select {
		case <-t.C:
			ctx := context.Background()

			if err := u.tryUpdateConfig(ctx); err != nil {
				log.Error().Err(err).Msgf("Can't update config")
			}

			if err := u.tryUpdateDetectors(ctx); err != nil {
				log.Error().Err(err).Msgf("Can't update detectors")
			}
		case <-stop:
			return
		}
	}
}

func (u *Updater) tryUpdateConfig(ctx context.Context) error {
	cfg, err := u.ConfigRepository.GetLast(ctx, true) // preload is on
	if err != nil {
		return fmt.Errorf("can't get last config from DB: %w", err)
	}

	oldCfg := u.Processor.Config()

	log.Debug().Interface("old_config", oldCfg).Msgf("Old processor config")
	log.Debug().Interface("new_config", cfg).Msgf("New processor config")

	sel, changed := config.Diff(oldCfg, cfg)
	if changed {
		log.Info().
			Interface("config", cfg).
			Interface("selector", sel).
			Msgf("Processor config changed")

		u.Processor.SetConfig(cfg)
	} else {
		log.Debug().Msgf("Processor config didn't change")
	}

	return nil
}

func (u *Updater) tryUpdateDetectors(ctx context.Context) error {
	hashes, err := u.DetectorRepository.GetAllHashes(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't get all hashes from DB: %w", err)
	}

	decoded, err := fromHex(hashes)
	if err != nil {
		return fmt.Errorf("can't decode hashes: %w", err)
	}

	rootHash := processor.HashesHashAsHex(decoded)
	_, oldRootHash := u.Processor.Bins()

	log.Debug().Interface("old_root_hash", oldRootHash).Msgf("Old detectors root hash")
	log.Debug().Interface("new_root_hash", rootHash).Msgf("New detectors root hash")

	if rootHash != oldRootHash {
		bins, err := u.DetectorRepository.GetAllBins(ctx, nil)
		if err != nil {
			return fmt.Errorf("can't get bins from DB: %w", err)
		}

		u.Processor.UpdateDetectors(bins)

		log.Info().Interface("root_hash", rootHash).Msgf("Detectors changed")
	} else {
		log.Debug().Msgf("Detectors didn't change")
	}

	return nil
}

func fromHex(hs []string) ([][]byte, error) {
	ds := make([][]byte, 0, len(hs))

	for _, h := range hs {
		d, err := hex.DecodeString(h)
		if err != nil {
			return nil, fmt.Errorf("can't decode '%s': %w", h, err)
		}

		ds = append(ds, d)
	}

	return ds, nil
}
