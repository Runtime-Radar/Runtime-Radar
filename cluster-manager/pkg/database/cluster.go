package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/cluster-manager/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClusterRepository interface {
	Add(ctx context.Context, cs ...*model.Cluster) error
	GetCount(ctx context.Context, filter any) (int, error)
	GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preloadData bool) ([]*model.Cluster, error)
	GetAll(ctx context.Context, filter, order any, preloadData bool) ([]*model.Cluster, error)
	GetByID(ctx context.Context, id uuid.UUID, preloadData bool) (*model.Cluster, error)
	GetByToken(ctx context.Context, token uuid.UUID, preloadData bool) (*model.Cluster, error)
	UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByFilter(ctx context.Context, filter any) error
}

type ClusterDatabase struct {
	*gorm.DB
}

func (cd *ClusterDatabase) Add(ctx context.Context, cs ...*model.Cluster) error {
	if len(cs) == 0 {
		return nil
	}

	return cd.WithContext(ctx).Create(cs).Error
}

func (cd *ClusterDatabase) GetPage(ctx context.Context, filter, order any, pageSize int, pageNum int, preloadData bool) ([]*model.Cluster, error) {
	cs := []*model.Cluster{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = cd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Limit(pageSize).
		Offset(pageSize * (pageNum - 1)).
		Find(&cs).
		Error

	return cs, err
}

func (cd *ClusterDatabase) GetAll(ctx context.Context, filter, order interface{}, preloadData bool) ([]*model.Cluster, error) {
	cs := []*model.Cluster{}

	if filter == nil {
		filter = ""
	}

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = cd.preloadData(ctx, preloadData).
		Where(filter).
		Order(sanitizedOrder).
		Find(&cs).
		Error

	return cs, err
}

func (cd *ClusterDatabase) GetCount(ctx context.Context, filter interface{}) (int, error) {
	if filter == nil {
		filter = ""
	}

	var count int64

	err := cd.WithContext(ctx).
		Model(&model.Cluster{}).
		Where(filter).
		Count(&count).
		Error

	return int(count), err
}

func (cd *ClusterDatabase) GetByID(ctx context.Context, id uuid.UUID, preloadData bool) (*model.Cluster, error) {
	if id == uuid.Nil {
		return nil, errors.New("incorrect cluster ID")
	}

	c := &model.Cluster{}

	err := cd.preloadData(ctx, preloadData).
		Unscoped(). // softly deleted rows should also be selected
		Where(&model.Cluster{Base: model.Base{ID: id}}).
		Take(&c).
		Error

	return c, err
}

func (cd *ClusterDatabase) GetByToken(ctx context.Context, token uuid.UUID, preloadData bool) (*model.Cluster, error) {
	if token == uuid.Nil {
		return nil, errors.New("incorrect token")
	}

	c := &model.Cluster{}

	err := cd.preloadData(ctx, preloadData).
		Where(&model.Cluster{Token: token}).
		Take(&c).
		Error

	return c, err
}

// UpdateWithMap updates record in DB by setting of provided key-value (even zeroed) map entries, where key can either be DB column, or struct field name.
func (cd *ClusterDatabase) UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]interface{}) error {
	res := cd.WithContext(ctx).
		Model(&model.Cluster{Base: model.Base{ID: id}}).
		Updates(m)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (cd *ClusterDatabase) Delete(ctx context.Context, id uuid.UUID) error {
	return cd.WithContext(ctx).
		Delete(&model.Cluster{Base: model.Base{ID: id}}).
		Error
}

func (cd *ClusterDatabase) DeleteByFilter(ctx context.Context, filter any) error {
	errEmptyFilter := errors.New("no filters passed")

	if filter == nil {
		return errEmptyFilter
	}
	if f, ok := filter.(string); ok && f == "" {
		return errEmptyFilter
	}

	return cd.WithContext(ctx).
		Where(filter).
		Delete(&model.Cluster{}).
		Error
}

func (cd *ClusterDatabase) preloadData(ctx context.Context, preloadData bool) *gorm.DB {
	if preloadData {
		return cd.WithContext(ctx).
			// This should load all associations without nested: https://gorm.io/docs/preload.html#Preload-All
			Preload(clause.Associations)
	}
	return cd.WithContext(ctx)
}
