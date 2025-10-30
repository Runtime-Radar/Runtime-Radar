package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"gorm.io/gorm"
)

// EventsCleaner provides a job that is responsible for
// limiting the number of records in runtime_events table in ClickHouse.
// When the table has more than Limit records, older events are deleted.
//
// Records are removed using lightweight delete mechanism,
// which means that it's not determined how soon the data will be deleted from disk.
// See https://clickhouse.com/docs/sql-reference/statements/delete for details.
type EventsCleaner struct {
	*gorm.DB

	Interval time.Duration
	Limit    int // limit of records to be stored in table
}

func (e *EventsCleaner) Run(stop <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stop
		cancel()
	}()

	log.Debug().Msg("Runtime events cleaner started")
	defer log.Debug().Msg("Runtime events cleaner stopped")

	t := time.NewTicker(e.Interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			start := time.Now()

			if err := e.clean(ctx); err != nil {
				log.Error().Msgf("Can't clean runtime events: %v", err)
			} else {
				log.Info().Stringer("delay", time.Since(start)).Msg("Runtime events cleaned")
			}

		case <-ctx.Done():
			return
		}
	}
}

// clean runs an SQL query which deletes old records using lightweight delete.
// The query may take a long time, so that no timeout is set up.
func (e *EventsCleaner) clean(ctx context.Context) error {
	var count int64

	err := e.WithContext(ctx).
		Model(&model.RuntimeEvent{}).
		Count(&count).
		Error
	if err != nil {
		return fmt.Errorf("can't count events: %w", err)
	}

	// Mutation is heavyweight process and may take a long time to proceed,
	// even though the limit may not be exceeded.
	if int(count) <= e.Limit {
		log.Debug().Int64("events_count", count).Msg("Runtime events limit is not exceeded")
		return nil
	}
	log.Debug().Int64("events_count", count).Msg("Runtime events limit is exceeded")

	const q = `
		delete from runtime_events
		where registered_at < (
			select registered_at
			from runtime_events
			order by registered_at desc
			limit 1 offset ?
		)
	`

	// As we use offset mechanism we substract 1
	// so that registered_at is compared to nth record, not to n+1'th one.
	if err := e.WithContext(ctx).Exec(q, e.Limit-1).Error; err != nil {
		return fmt.Errorf("can't delete events: %w", err)
	}

	return nil
}
