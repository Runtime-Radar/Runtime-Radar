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

type Syslog struct {
	Base
	Name          string `gorm:"index"`
	Address       string
	Notifications []*Notification `gorm:"polymorphic:Integration;polymorphicValue:syslog"`
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
	// Meta is not stored at database, but added at runtime
	Meta IntegrationMeta `gorm:"-"`
}

func (s *Syslog) BeforeCreate(tx *gorm.DB) error {
	base := any(&s.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if s.Name != "" {
		if err := s.checkNameUnique(tx, s.Name); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syslog) BeforeUpdate(tx *gorm.DB) error {
	base := any(&s.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := s.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syslog) AfterDelete(tx *gorm.DB) error {
	base := any(&s.Base)
	if b, ok := base.(callbacks.AfterDeleteInterface); ok {
		if err := b.AfterDelete(tx); err != nil {
			return err
		}
	}

	return tx.Where(&Notification{IntegrationID: s.ID}).
		Delete(&Notification{}).
		Error
}

func (s *Syslog) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Syslog{Name: name}).
		Not(&Syslog{Base: Base{ID: s.ID}}). // this will be needed in future if updates are implemented
		Take(&Syslog{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrIntegrationNameInUse
}

func (s *Syslog) GetID() uuid.UUID {
	return s.ID
}

func (s *Syslog) SetMeta(m IntegrationMeta) {
	s.Meta = m
}

func (s *Syslog) EncryptSensitive(_ cipher.Crypter) {
}

func (s *Syslog) DecryptSensitive(_ cipher.Crypter) {
}

func (s *Syslog) MaskSensitive() {
}
