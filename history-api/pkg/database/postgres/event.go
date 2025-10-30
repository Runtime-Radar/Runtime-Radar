// nolint: goconst
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Add(ctx context.Context, ss ...*model.Event) error
	GetByID(ctx context.Context, id uuid.UUID, preload *EventPreloadConfig) (*model.Event, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID, order any, preload *EventPreloadConfig) ([]*model.Event, error)
	Update(ctx context.Context, s *model.Event) error
	UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]any) error
	Save(ctx context.Context, s *model.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preload *EventPreloadConfig) ([]*model.Event, error)
	GetCount(ctx context.Context, filter any) (int, error)
	GetAll(ctx context.Context, filter, order any, preload *EventPreloadConfig) ([]*model.Event, error)
}

type EventDatabase struct {
	*gorm.DB
}

type EventPreloadConfig struct {
	WithIncident bool
}

// Add adds new entry to the database, it can add multiple instances at once.
func (ed *EventDatabase) Add(ctx context.Context, es ...*model.Event) error {
	if len(es) == 0 {
		return nil
	}

	return ed.WithContext(ctx).Create(es).Error
}

func (ed *EventDatabase) GetByID(ctx context.Context, id uuid.UUID, preload *EventPreloadConfig) (*model.Event, error) {
	if id == uuid.Nil {
		return nil, errors.New("incorrect event ID")
	}

	e := &model.Event{}

	err := ed.preloadData(ctx, preload).
		Where(&model.Event{Base: model.Base{ID: id}}).
		Take(&e).
		Error

	return e, err
}

func (ed *EventDatabase) GetByIDs(ctx context.Context, ids []uuid.UUID, order any, preload *EventPreloadConfig) ([]*model.Event, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	es := make([]*model.Event, 0, len(ids))

	err := ed.preloadData(ctx, preload).
		Where(ids).
		Order(order).
		Find(&es).
		Error

	return es, err
}

// Update updates record in DB by setting of non-zero fields of given struct.
func (ed *EventDatabase) Update(ctx context.Context, e *model.Event) error {
	return ed.WithContext(ctx).
		Updates(e).
		Error
}

// UpdateWithMap updates record in DB by setting of provided key-value (even zeroed) map entries, where key can either be DB column, or struct field name.
func (ed *EventDatabase) UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]any) error {
	return ed.WithContext(ctx).
		Model(&model.Event{Base: model.Base{ID: id}}).
		Updates(m).
		Error
}

// Save updates record in DB by setting of ALL (even zeroed) fields of given struct,
// it will INSERT a record if primary key (ID) not set.
func (ed *EventDatabase) Save(ctx context.Context, e *model.Event) error {
	return ed.WithContext(ctx).
		Save(e).
		Error
}

func (ed *EventDatabase) Delete(ctx context.Context, id uuid.UUID) error {
	return ed.WithContext(ctx).
		Delete(&model.Event{Base: model.Base{ID: id}}).
		Error
}

func (ed *EventDatabase) DeleteByFilter(ctx context.Context, filter any) error {
	errEmptyFilter := errors.New("no filters passed")

	if filter == nil {
		return errEmptyFilter
	}
	if f, ok := filter.(string); ok && f == "" {
		return errEmptyFilter
	}

	return ed.WithContext(ctx).
		Where(filter).
		Delete(&model.Event{}).
		Error
}

func (ed *EventDatabase) GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preload *EventPreloadConfig) ([]*model.Event, error) {
	es := []*model.Event{}

	if filter == nil {
		filter = ""
	}

	err := ed.preloadData(ctx, preload).
		Where(filter).
		Order(order).
		Limit(pageSize).
		Offset(pageSize * (pageNum - 1)).
		Find(&es).
		Error

	return es, err
}

func (ed *EventDatabase) GetCount(ctx context.Context, filter any) (int, error) {
	if filter == nil {
		filter = ""
	}

	var count int64

	err := ed.WithContext(ctx).
		Model(&model.Event{}).
		Where(filter).
		Count(&count).
		Error

	return int(count), err
}

func (ed *EventDatabase) GetAll(ctx context.Context, filter, order any, preload *EventPreloadConfig) ([]*model.Event, error) {
	es := []*model.Event{}

	if filter == nil {
		filter = ""
	}

	err := ed.preloadData(ctx, preload).
		Where(filter).
		Order(order).
		Find(&es).
		Error

	return es, err
}

func (ed *EventDatabase) preloadData(ctx context.Context, preload *EventPreloadConfig) *gorm.DB {
	db := ed.WithContext(ctx)

	if preload == nil {
		return db
	}
	if preload.WithIncident {
		db = db.Preload("Incident")
	}

	return db
}
