package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/auth-center/pkg/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	GetAll(ctx context.Context) ([]*model.Role, error)
}

type RoleDatabase struct {
	*gorm.DB
}

func (rd *RoleDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	if id == uuid.Nil {
		return nil, errors.New("empty id")
	}

	role := &model.Role{}

	err := rd.WithContext(ctx).
		Where(&model.Role{ID: id}).
		Take(&role).Error

	return role, err
}

func (rd *RoleDatabase) GetAll(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role

	err := rd.WithContext(ctx).Find(&roles).Error

	return roles, err
}
