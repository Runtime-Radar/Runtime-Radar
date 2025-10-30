package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleRepository interface {
	Add(ctx context.Context, ss ...*model.Rule) error
	GetByID(ctx context.Context, id uuid.UUID, preloadData bool) (*model.Rule, error)
	GetByTypeAndIDs(ctx context.Context, rt model.RuleType, ids []uuid.UUID, order interface{}, preloadData bool) ([]*model.Rule, error)
	Update(ctx context.Context, s *model.Rule) error
	UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]interface{}) error
	Save(ctx context.Context, s *model.Rule) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPage(ctx context.Context, filter, order interface{}, pageSize int, pageNum int, preloadData bool) ([]*model.Rule, error)
	GetCount(ctx context.Context, filter interface{}) (int, error)
	GetAll(ctx context.Context, filter, order interface{}, preloadData bool) ([]*model.Rule, error)
}

type RuleDatabase struct {
	*gorm.DB
}

// Add adds new entry to the database, it can add multiple instances at once.
func (rd *RuleDatabase) Add(ctx context.Context, rs ...*model.Rule) error {
	if len(rs) == 0 {
		return nil
	}

	return rd.WithContext(ctx).Create(rs).Error
}

// GetByID returns rule by id. This method allows getting record that was softly deleted as well.
func (rd *RuleDatabase) GetByID(ctx context.Context, id uuid.UUID, preloadData bool) (*model.Rule, error) {
	if id == uuid.Nil {
		return nil, errors.New("incorrect rule ID")
	}

	r := &model.Rule{}

	err := rd.preloadData(ctx, preloadData).
		Unscoped(). // softly deleted rows should also be selected
		Where(&model.Rule{Base: model.Base{ID: id}}).
		Take(&r).
		Error

	return r, err
}

func (rd *RuleDatabase) GetByTypeAndIDs(ctx context.Context, rt model.RuleType, ids []uuid.UUID, order interface{}, preloadData bool) ([]*model.Rule, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	rs := make([]*model.Rule, 0, len(ids))

	err = rd.preloadData(ctx, preloadData).
		Where(ids).
		Where(&model.Rule{Type: rt}).
		Order(sanitizedOrder).
		Find(&rs).
		Error

	return rs, err
}

// Update updates record in DB by setting of non-zero fields of given struct.
func (rd *RuleDatabase) Update(ctx context.Context, r *model.Rule) error {
	return rd.WithContext(ctx).
		Updates(r).
		Error
}

// UpdateWithMap updates record in DB by setting of provided key-value (even zeroed) map entries, where key can either be DB column, or struct field name.
func (rd *RuleDatabase) UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]interface{}) error {
	return rd.WithContext(ctx).
		Model(&model.Rule{Base: model.Base{ID: id}}).
		Updates(m).
		Error
}

// Save updates record in DB by setting of ALL (even zeroed) fields of given struct,
// it will INSERT a record if primary key (ID) not set.
func (rd *RuleDatabase) Save(ctx context.Context, r *model.Rule) error {
	return rd.WithContext(ctx).
		Save(r).
		Error
}

func (rd *RuleDatabase) Delete(ctx context.Context, id uuid.UUID) error {
	return rd.WithContext(ctx).
		Delete(&model.Rule{Base: model.Base{ID: id}}). // model.Rule has DeletedAt field => soft delete will happen
		Error
}

func (rd *RuleDatabase) GetPage(ctx context.Context, filter, order interface{}, pageSize int, pageNum int, preloadData bool) ([]*model.Rule, error) {
	rs := []*model.Rule{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = rd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Limit(pageSize).
		Offset(pageSize * (pageNum - 1)).
		Find(&rs).
		Error

	return rs, err
}

func (rd *RuleDatabase) GetCount(ctx context.Context, filter interface{}) (int, error) {
	if filter == nil {
		filter = ""
	}

	var count int64

	err := rd.WithContext(ctx).
		Model(&model.Rule{}).
		Where(filter).
		Count(&count).
		Error

	return int(count), err
}

func (rd *RuleDatabase) GetAll(ctx context.Context, filter, order interface{}, preloadData bool) ([]*model.Rule, error) {
	rs := []*model.Rule{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = rd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Find(&rs).
		Error

	return rs, err
}

func (rd *RuleDatabase) preloadData(ctx context.Context, preloadData bool) *gorm.DB {
	if preloadData {
		return rd.WithContext(ctx).
			// This should load all associations without nested: https://gorm.io/docs/preload.html#Preload-All
			Preload(clause.Associations)
	}
	return rd.WithContext(ctx)
}
