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

type Webhook struct {
	Base
	Name string `gorm:"index"`
	URL  string
	// Login and Password represent credentials for basic auth
	Login             string
	Password          string `gorm:"-"`
	EncryptedPassword string
	Insecure          bool
	CA                string
	Notifications     []*Notification `gorm:"polymorphic:Integration;polymorphicValue:email"`
	DeletedAt         gorm.DeletedAt  `gorm:"index"`
	// Meta is not stored at database, but added at runtime
	Meta IntegrationMeta `gorm:"-"`
}

func (w *Webhook) BeforeCreate(tx *gorm.DB) error {
	base := any(&w.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if w.Name != "" {
		if err := w.checkNameUnique(tx, w.Name); err != nil {
			return err
		}
	}

	return nil
}

func (w *Webhook) BeforeUpdate(tx *gorm.DB) error {
	base := any(&w.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := w.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (w *Webhook) AfterDelete(tx *gorm.DB) error {
	base := any(&w.Base)
	if b, ok := base.(callbacks.AfterDeleteInterface); ok {
		if err := b.AfterDelete(tx); err != nil {
			return err
		}
	}

	return tx.Where(&Notification{IntegrationID: w.ID}).
		Delete(&Notification{}).
		Error
}

func (w *Webhook) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Webhook{Name: name}).
		Not(&Webhook{Base: Base{ID: w.ID}}). // this will be needed in future if updates are implemented
		Take(&Webhook{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrIntegrationNameInUse
}

func (w *Webhook) GetID() uuid.UUID {
	return w.ID
}

func (w *Webhook) EncryptSensitive(c cipher.Crypter) {
	// nothing to encrypt
	if w.Password == "" {
		return
	}
	w.EncryptedPassword = c.EncryptStringAsHex(w.Password)
}

func (w *Webhook) DecryptSensitive(c cipher.Crypter) {
	// we've either have nothing to decrypt or password is already decrypted
	if w.EncryptedPassword == "" || w.Password != "" {
		return
	}
	w.Password = c.DecryptHexAsString(w.EncryptedPassword)
}

func (w *Webhook) MaskSensitive() {
	const mask = "********"
	w.EncryptedPassword = mask
	w.Password = mask
}

func (w *Webhook) SetMeta(m IntegrationMeta) {
	w.Meta = m
}
