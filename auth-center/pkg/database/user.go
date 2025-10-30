package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/runtime-radar/runtime-radar/auth-center/pkg/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrUsernameAlreadyExists = errors.New("username already exists")

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	GetUsersByRoleID(ctx context.Context, roleID uuid.UUID) ([]*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, pass string, ID uuid.UUID) error
}

type UserDatabase struct {
	*gorm.DB
}

func (ud *UserDatabase) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("empty id")
	}

	user := &model.User{}

	err := ud.WithContext(ctx).
		Preload("Role", nil).
		Where(&model.User{
			Base: model.Base{ID: id}}).
		First(user).Error

	return user, err
}

func (ud *UserDatabase) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("empty username")
	}

	user := &model.User{}

	err := ud.WithContext(ctx).
		Preload("Role", nil).
		Where(model.User{
			Username: username}).
		First(user).Error

	return user, err
}

func (ud *UserDatabase) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	err := ud.WithContext(ctx).
		Preload("Role", nil).
		Find(&users).Error

	return users, err
}

func (ud *UserDatabase) GetUsersByRoleID(ctx context.Context, roleID uuid.UUID) ([]*model.User, error) {
	var users []*model.User

	err := ud.WithContext(ctx).
		Preload("Role", nil).
		Where(model.User{
			RoleID: roleID}).
		Find(&users).Error

	return users, err
}

func (ud *UserDatabase) Create(ctx context.Context, user *model.User) error {
	err := ud.WithContext(ctx).Create(user).Error

	return ud.validatePGError(err)
}

func (ud *UserDatabase) Update(ctx context.Context, user *model.User) error {
	err := ud.WithContext(ctx).
		Model(&model.User{
			Base: model.Base{ID: user.ID}}).
		Updates(&model.User{
			Email:         user.Email,
			RoleID:        user.RoleID,
			MappingRoleID: user.MappingRoleID,
		}).Error

	return ud.validatePGError(err)
}

func (ud *UserDatabase) Delete(ctx context.Context, id uuid.UUID) error {
	query := ud.WithContext(ctx).
		Delete(&model.User{
			Base: model.Base{ID: id},
		})
	err := query.Error

	if err != nil {
		return err
	}

	if query.RowsAffected < 1 {
		return fmt.Errorf("row with id='%s' cannot be deleted because it doesn't exist", id.String())
	}

	return nil
}

func (ud *UserDatabase) ChangePassword(ctx context.Context, pass string, id uuid.UUID) error {
	return ud.WithContext(ctx).
		Model(&model.User{
			Base: model.Base{ID: id}}).
		Updates(&model.User{
			HashedPassword:        pass,
			LastPasswordChangedAt: time.Now()},
		).Error
}

func (ud *UserDatabase) preloadData(ctx context.Context, preloadData bool) *gorm.DB {
	if preloadData {
		return ud.WithContext(ctx).
			// This should load all associations without nested: https://gorm.io/docs/preload.html#Preload-All
			Preload(clause.Associations)
	}
	return ud.WithContext(ctx)
}

func (ud *UserDatabase) validatePGError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return gorm.ErrForeignKeyViolated
		}

		if pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "uni_users_username":
				return ErrUsernameAlreadyExists
			default:
				return gorm.ErrCheckConstraintViolated
			}
		}
	}

	return err
}
