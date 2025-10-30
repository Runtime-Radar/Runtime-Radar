package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"gorm.io/gorm"
)

// runtimeEventColumns is the set of columns to be selected from runtime_events table.
// In fact, table contains many more columns, but they're mainly used for filtering.
const runtimeEventColumns = "id,source_event,threats,registered_at,tetragon_version,is_incident,incident_severity,block_by,notify_by,detect_errors"

type RuntimeEventRepository interface {
	Add(ctx context.Context, events *[]model.RuntimeEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.RuntimeEvent, error)
	GetRightSlice(ctx context.Context, cursor time.Time, filter any, sliceSize int) ([]*model.RuntimeEvent, error)
	GetLeftSlice(ctx context.Context, cursor time.Time, filter any, sliceSize int) ([]*model.RuntimeEvent, error)
}

type RuntimeEventDatabase struct {
	*gorm.DB
}

func (d *RuntimeEventDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.RuntimeEvent, error) {
	if id == uuid.Nil {
		return nil, errors.New("incorrect runtime event ID")
	}

	e := &model.RuntimeEvent{}

	err := d.WithContext(ctx).
		Select(runtimeEventColumns).
		Where(&model.RuntimeEvent{ID: id}).
		Take(&e).
		Error

	return e, err
}

func (d *RuntimeEventDatabase) GetRightSlice(ctx context.Context, cursor time.Time, filter any, sliceSize int) ([]*model.RuntimeEvent, error) {
	var events []*model.RuntimeEvent

	if filter == nil {
		filter = ""
	}

	err := d.WithContext(ctx).
		Select(runtimeEventColumns).
		Where(filter).
		// toDateTime64 has to be used because DateTime64 cannot be automatically converted from string. See https://clickhouse.com/docs/en/sql-reference/data-types/datetime64 for details.
		Where("registered_at < toDateTime64(?, 9, ?)", cursor.Format(DateTimeFormat), cursor.Location().String()).
		Order("registered_at DESC").
		Limit(sliceSize).
		Find(&events).
		Error

	return events, err
}

func (d *RuntimeEventDatabase) GetLeftSlice(ctx context.Context, cursor time.Time, filter any, sliceSize int) ([]*model.RuntimeEvent, error) {
	var events []*model.RuntimeEvent

	if filter == nil {
		filter = ""
	}

	err := d.WithContext(ctx).
		Select(runtimeEventColumns).
		Where(filter).
		// toDateTime64 has to be used because DateTime64 cannot be automatically converted from string. See https://clickhouse.com/docs/en/sql-reference/data-types/datetime64 for details.
		Where("registered_at > toDateTime64(?, 9, ?)", cursor.Format(DateTimeFormat), cursor.Location().String()).
		Order("registered_at ASC").
		Limit(sliceSize).
		Find(&events).
		Error

	slices.Reverse(events)

	return events, err
}

type RuntimeEventBatchingDatabase struct {
	events    []model.RuntimeEvent
	batchSize int
	ticker    *time.Ticker

	*RuntimeEventDatabase
}

func NewRuntimeEventBatchingDatabase(batchSize int, flushInterval time.Duration, db *gorm.DB) *RuntimeEventBatchingDatabase {
	bdb := &RuntimeEventBatchingDatabase{
		make([]model.RuntimeEvent, 0, batchSize),
		batchSize,
		time.NewTicker(flushInterval),

		&RuntimeEventDatabase{db},
	}

	return bdb
}

// Add wraps RuntimeEventDatabase.Add by adding a batching mechanism:
//   - if number of current and buffered events reaches batchSize, they are flushed to an underlying DB,
//   - if ticker fires, current and buffered events also flushed to an underlying DB,
//   - otherwise current event(s) are added to internal slice for next iteration.
//
// In order to keep rate under incresing load it's possible to rise up batchSize to very high numbers
// such as 1e6 or even higher, thus trading off process memory for securing the rate.
//
// If for some reason after collecting some events, they stopped arriving, buffer won't be
// flushed until there is another Add call, which is similar to how most of systems with buffering work.
// However, this is very unlikely as events flow will be high most of the time, and if it's not, average
// events rate is expected to be more or less persistent, so it won't stop forever, but it may require
// few seconds or minutes depending on system configuration, such as enabled policies and filters.
//
// events is given as pointer to slice of models (unlike slice of pointers as it's done in most other places)
// for two reasons:
//   - weaken GC pressure due to allocations, as we effectively allocating a single backing array once when
//     RuntimeEventBatchingDatabase constructor initializes slice, then reusing memory (model consists of a fields of simple types)
//   - make GORM fill-in primary keys and/or let hooks reflect changes if they are, some details:
//     https://gorm.io/docs/create.html#Batch-Insert
//
// Memory consumption, including allocations, should be checked under load, when we have more time for that.
//
// NOTE: this implementation is kept UNSAFE for concurrent use intentionally, because:
//   - events are consumed in single goroutine, and it's highly likely to remain in future,
//   - avoiding mutexes and channels eliminates complex synchronization issues, and
//   - eliminates performance bottlenecks due to mutual blocking and more memory copying
//
// If we decide to go for making it safe for concurrent use later on, we should take a look at possibility to
// instantiate separate RuntimeEventBatchingDatabase for each of the goroutine first.
func (d *RuntimeEventBatchingDatabase) Add(ctx context.Context, events *[]model.RuntimeEvent) error {
	if events == nil {
		return errors.New("nil pointer to events slice")
	}

	select {
	case <-d.ticker.C:
		d.events = append(d.events, *events...)
		if err := d.RuntimeEventDatabase.Add(ctx, &d.events); err != nil {
			return fmt.Errorf("can't store runtime events: %w", err)
		}

		d.events = d.events[:0]
	default:
		if len(d.events)+len(*events) < d.batchSize {
			d.events = append(d.events, *events...)
		} else {
			d.events = append(d.events, *events...)
			if err := d.RuntimeEventDatabase.Add(ctx, &d.events); err != nil {
				return fmt.Errorf("can't store runtime events: %w", err)
			}
			log.Warn().Msgf("Events batch has been stored with higher rate than expected")

			d.events = d.events[:0]
		}
	}

	return nil
}

func (d *RuntimeEventDatabase) Add(ctx context.Context, events *[]model.RuntimeEvent) error {
	if events == nil {
		return errors.New("nil pointer to events slice")
	}

	if len(*events) == 0 {
		return nil
	}

	return d.WithContext(ctx).Create(events).Error
}
