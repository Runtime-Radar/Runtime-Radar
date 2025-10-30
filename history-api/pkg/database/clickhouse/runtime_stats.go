package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"gorm.io/gorm"
)

type StatsRepository interface {
	CountEvents(ctx context.Context, from, to time.Time, filter any) (int, error)
}

type StatsDatabase struct {
	*gorm.DB
}

func (s *StatsDatabase) CountEvents(ctx context.Context, from, to time.Time, filter any) (int, error) {
	if reason, ok := validateTimePeriod(from, to); !ok {
		return 0, fmt.Errorf("invalid time period: %s", reason)
	}

	if filter == nil {
		filter = ""
	}

	qb := s.WithContext(ctx).
		Model(&model.RuntimeEvent{}).
		Where(filter)

	qb = applyTimePeriod(qb, "registered_at", from, to)

	var cnt int64
	err := qb.Count(&cnt).Error

	return int(cnt), err
}

func validateTimePeriod(from, to time.Time) (reason string, ok bool) {
	if from.IsZero() {
		return "from is not set", false
	}
	if to.IsZero() {
		return "to is not set", false
	}
	if from.After(to) {
		return "from is after to", false
	}
	return "", true
}

// applyTimePeriod returns db with time period added to it. It only applies from/to if it's not zero-valued.
func applyTimePeriod(db *gorm.DB, field string, from, to time.Time) *gorm.DB {
	if !from.IsZero() {
		expr := fmt.Sprintf("%s > toDateTime64(?, 9, ?)", field)
		// time is converted to UTC explicitly, because time.Now() returns timestamp with "Local" location unless timezone is configured explicitly.
		// "Local" is not accepted by Clickhouse.
		db = db.Where(expr, from.UTC().Format(DateTimeFormat), from.UTC().Location().String())
	}

	if !to.IsZero() {
		expr := fmt.Sprintf("%s < toDateTime64(?, 9, ?)", field)
		// time is converted to UTC explicitly, because time.Now() returns timestamp with "Local" location unless timezone is configured explicitly.
		// "Local" is not accepted by Clickhouse.
		db = db.Where(expr, to.UTC().Format(DateTimeFormat), to.UTC().Location().String())
	}

	return db
}
