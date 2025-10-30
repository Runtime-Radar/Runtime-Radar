package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"gorm.io/gorm"
)

type IntegrationRepository interface {
	Add(context.Context, model.Integration) error
	GetByTypeAndID(ctx context.Context, it string, id uuid.UUID) (model.Integration, error)
	GetAllByType(ctx context.Context, it string, order any) ([]model.Integration, error)
	GetByNotifications(context.Context, ...*model.Notification) ([]model.Integration, error)
	DeleteByTypeAndID(context.Context, string, uuid.UUID) error
	UpdateWithMap(ctx context.Context, it string, id uuid.UUID, m map[string]any) error
}

type IntegrationDatabase struct {
	*gorm.DB
}

func (i *IntegrationDatabase) Add(ctx context.Context, integration model.Integration) error {
	return i.WithContext(ctx).
		Create(integration.(any)).
		Error
}

func (i *IntegrationDatabase) GetByTypeAndID(ctx context.Context, it string, id uuid.UUID) (model.Integration, error) {
	integration, ok := integrationFromTypeAndID(it, id)
	if !ok {
		return nil, fmt.Errorf("invalid integration type: %s", it)
	}

	err := i.WithContext(ctx).
		Take(integration.(any)).
		Error

	return integration, err
}

func (i *IntegrationDatabase) DeleteByTypeAndID(ctx context.Context, it string, id uuid.UUID) error {
	integration, ok := integrationFromTypeAndID(it, id)
	if !ok {
		return fmt.Errorf("invalid integration type: %s", it)
	}

	return i.WithContext(ctx).
		Delete(integration.(any)).
		Error
}

func (i *IntegrationDatabase) GetAllByType(ctx context.Context, it string, order any) ([]model.Integration, error) {
	sanitizedOrder, err := sanitizeOrder(order)
	if err != nil {
		return nil, err
	}

	switch it {
	case model.IntegrationEmail:
		es, err := i.getEmails(ctx, nil, sanitizedOrder)
		if err != nil {
			return nil, fmt.Errorf("can't get emails: %w", err)
		}

		return toIntegrationSlice(es), nil

	case model.IntegrationWebhook:
		ws, err := i.getWebhooks(ctx, nil, sanitizedOrder)
		if err != nil {
			return nil, fmt.Errorf("can't get webhooks: %w", err)
		}

		return toIntegrationSlice(ws), nil
	case model.IntegrationSyslog:
		sl, err := i.getSyslogs(ctx, nil, sanitizedOrder)
		if err != nil {
			return nil, fmt.Errorf("can't get syslogs: %w", err)
		}

		return toIntegrationSlice(sl), nil
	default:
		return nil, fmt.Errorf("invalid integration type: %s", it)
	}
}

func (i *IntegrationDatabase) UpdateWithMap(ctx context.Context, it string, id uuid.UUID, m map[string]any) error {
	integration, ok := integrationFromTypeAndID(it, id)
	if !ok {
		return fmt.Errorf("invalid integration type: %s", it)
	}

	return i.WithContext(ctx).
		Model(integration.(any)).
		Updates(m).
		Error
}

func (i *IntegrationDatabase) GetByNotifications(ctx context.Context, ns ...*model.Notification) ([]model.Integration, error) {
	idsByType := make(map[string][]uuid.UUID)
	for _, n := range ns {
		idsByType[n.IntegrationType] = append(idsByType[n.IntegrationType], n.IntegrationID)
	}

	var res []model.Integration

	for it, ids := range idsByType {
		switch it {
		case model.IntegrationEmail:
			es, err := i.getEmails(ctx, ids, nil)
			if err != nil {
				return nil, fmt.Errorf("can't get emails: %w", err)
			}

			res = append(res, toIntegrationSlice(es)...)

		case model.IntegrationWebhook:
			ws, err := i.getWebhooks(ctx, ids, nil)
			if err != nil {
				return nil, fmt.Errorf("can't get webhooks: %w", err)
			}

			res = append(res, toIntegrationSlice(ws)...)
		case model.IntegrationSyslog:
			sl, err := i.getSyslogs(ctx, ids, nil)
			if err != nil {
				return nil, fmt.Errorf("can't get syslogs: %w", err)
			}

			res = append(res, toIntegrationSlice(sl)...)
		}
	}

	return res, nil
}

func (i *IntegrationDatabase) getEmails(ctx context.Context, filter any, order any) ([]*model.Email, error) {
	var es []*model.Email

	if filter == nil {
		filter = ""
	}

	err := i.WithContext(ctx).
		Where(filter).
		Order(order).
		Find(&es).
		Error

	return es, err
}

func (i *IntegrationDatabase) getWebhooks(ctx context.Context, filter, order any) ([]*model.Webhook, error) {
	var ws []*model.Webhook

	if filter == nil {
		filter = ""
	}

	err := i.WithContext(ctx).
		Where(filter).
		Order(order).
		Find(&ws).
		Error

	return ws, err
}

func (i *IntegrationDatabase) getSyslogs(ctx context.Context, filter, order any) ([]*model.Syslog, error) {
	var sl []*model.Syslog

	if filter == nil {
		filter = ""
	}

	err := i.WithContext(ctx).
		Where(filter).
		Order(order).
		Find(&sl).
		Error

	return sl, err
}

// IntegrationFromTypeAndID returns integration gorm model with ID set depending on given integration type
func integrationFromTypeAndID(it string, id uuid.UUID) (i model.Integration, ok bool) {
	switch it {
	case model.IntegrationEmail:
		return &model.Email{Base: model.Base{ID: id}}, true
	case model.IntegrationWebhook:
		return &model.Webhook{Base: model.Base{ID: id}}, true
	case model.IntegrationSyslog:
		return &model.Syslog{Base: model.Base{ID: id}}, true
	default:
		return nil, false
	}
}

func toIntegrationSlice[T model.Integration](impls []T) []model.Integration {
	res := make([]model.Integration, 0, len(impls))
	for _, impl := range impls {
		res = append(res, impl)
	}
	return res
}
