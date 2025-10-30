package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

var (
	ErrTokenNameInUse = errors.New("token name in use")
)

type Permissions jwt.RolePermissions

type AccessToken struct {
	Base

	Name          string
	UserID        uuid.UUID
	Hash          string       `gorm:"uniqueIndex"`
	Permissions   *Permissions `gorm:"type:jsonb"`
	ExpiresAt     *time.Time
	InvalidatedAt *time.Time // Timestamp when the token was invalidated by an administrator. Users cannot invalidate their own tokens.
}

func (a *AccessToken) BeforeCreate(tx *gorm.DB) error {
	base := any(&a.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if err := a.checkTokenNameUnique(tx); err != nil {
		return err
	}

	return nil
}

func (a *AccessToken) checkTokenNameUnique(tx *gorm.DB) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&AccessToken{Name: a.Name, UserID: a.UserID}).
		Not(&AccessToken{Base: Base{ID: a.ID}}).
		Take(&AccessToken{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if token name is in use: %w", err)
	}

	return ErrTokenNameInUse
}

func (p *Permissions) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, p)
}

func (p *Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}
