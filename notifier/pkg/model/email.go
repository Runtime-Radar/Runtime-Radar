package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/security/cipher"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type EmailAuthType byte

const (
	AuthTypePlain EmailAuthType = iota
	AuthTypeLogin
	AuthTypeCramMD5
	AuthTypeNone
)

type Email struct {
	Base
	Name              string `gorm:"index"`
	From              string
	Server            string
	AuthType          EmailAuthType
	Username          string
	EncryptedPassword string
	Password          string `gorm:"-"`
	UseTLS            bool   `gorm:"not null;default:false"`
	UseStartTLS       bool   `gorm:"not null;default:false"`
	Insecure          bool   `gorm:"not null;default:false"`
	CA                string
	Notifications     []*Notification `gorm:"polymorphic:Integration;polymorphicValue:email"`
	DeletedAt         gorm.DeletedAt  `gorm:"index"`
	// Meta is not stored at database, but added at runtime
	Meta IntegrationMeta `gorm:"-"`
}

func (e *Email) BeforeCreate(tx *gorm.DB) error {
	base := any(&e.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if e.Name != "" {
		if err := e.checkNameUnique(tx, e.Name); err != nil {
			return err
		}
	}

	return nil
}

func (e *Email) BeforeUpdate(tx *gorm.DB) error {
	base := any(&e.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := e.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (e *Email) AfterDelete(tx *gorm.DB) error {
	base := any(&e.Base)
	if b, ok := base.(callbacks.AfterDeleteInterface); ok {
		if err := b.AfterDelete(tx); err != nil {
			return err
		}
	}

	return tx.Where(&Notification{IntegrationID: e.ID}).
		Delete(&Notification{}).
		Error
}

func (e *Email) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Email{Name: name}).
		Not(&Email{Base: Base{ID: e.ID}}). // this will be needed in future if updates are implemented
		Take(&Email{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrIntegrationNameInUse
}

func (e *Email) GetID() uuid.UUID {
	return e.ID
}

func (e *Email) EncryptSensitive(c cipher.Crypter) {
	// nothing to encrypt
	if e.Password == "" {
		return
	}
	e.EncryptedPassword = c.EncryptStringAsHex(e.Password)
}

func (e *Email) DecryptSensitive(c cipher.Crypter) {
	// we either have nothing to decrypt or password is already decrypted
	if e.EncryptedPassword == "" || e.Password != "" {
		return
	}
	e.Password = c.DecryptHexAsString(e.EncryptedPassword)
}

func (e *Email) MaskSensitive() {
	const mask = "********"
	e.EncryptedPassword = mask
	e.Password = mask
}

func (e *Email) SetMeta(m IntegrationMeta) {
	e.Meta = m
}
