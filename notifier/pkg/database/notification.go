package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Add(ctx context.Context, ns ...*model.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID, order any) ([]*model.Notification, error)
	GetAll(ctx context.Context, filter, order any) ([]*model.Notification, error)
	GetByIntegrationID(ctx context.Context, integrationID uuid.UUID, order any) ([]*model.Notification, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]any) error
}

type NotificationDatabase struct {
	*gorm.DB
}

var _ NotificationRepository = (*NotificationDatabase)(nil)

func (nd *NotificationDatabase) Add(ctx context.Context, ns ...*model.Notification) error {
	if len(ns) == 0 {
		return nil
	}

	return nd.WithContext(ctx).Create(ns).Error
}

// GetByID returns notification by id. This method allows getting record that was softly deleted as well.
func (nd *NotificationDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	if id == uuid.Nil {
		return nil, errors.New("incorrect notification ID")
	}

	n := &model.Notification{}

	err := nd.
		WithContext(ctx).
		Unscoped(). // softly deleted rows should also be selected
		Where(&model.Notification{Base: model.Base{ID: id}}).
		Take(&n).
		Error

	return n, err
}

func (nd *NotificationDatabase) GetByIDs(ctx context.Context, ids []uuid.UUID, order any) ([]*model.Notification, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	ns := make([]*model.Notification, 0, len(ids))

	err = nd.WithContext(ctx).
		Where(ids).
		Order(sanitizedOrder).
		Find(&ns).
		Error

	return ns, err
}

func (nd *NotificationDatabase) GetAll(ctx context.Context, filter, order any) ([]*model.Notification, error) {
	ns := make([]*model.Notification, 0)

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = ""
	}

	err = nd.
		WithContext(ctx).
		Where(filter).
		Order(sanitizedOrder).
		Find(&ns).
		Error

	return ns, err
}

func (nd *NotificationDatabase) GetByIntegrationID(ctx context.Context, integrationID uuid.UUID, order any) ([]*model.Notification, error) {
	ns := make([]*model.Notification, 0)

	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	err = nd.
		WithContext(ctx).
		Where(&model.Notification{IntegrationID: integrationID}).
		Order(sanitizedOrder).
		Find(&ns).
		Error

	return ns, err
}

func (nd *NotificationDatabase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return nd.WithContext(ctx).
		Delete(&model.Notification{Base: model.Base{ID: id}}).
		Error
}

func (nd *NotificationDatabase) UpdateWithMap(ctx context.Context, id uuid.UUID, m map[string]any) error {
	return nd.WithContext(ctx).
		Model(&model.Notification{Base: model.Base{ID: id}}).
		Updates(m).
		Error
}
