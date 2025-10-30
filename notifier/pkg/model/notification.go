package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/runtime-radar/runtime-radar/notifier/api"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type (
	EmailConfig   api.EmailConfig
	WebhookConfig api.WebhookConfig
	SyslogConfig  api.SyslogConfig
)

// Notification is used to build final message that should be sent via some transport (IntegrationID).
// Only one config (EmailConfig, WebhookConfig, etc.) should be filled at a time depending on IntegrationType.
type Notification struct {
	Base
	Name            string         `gorm:"index"`
	Recipients      pq.StringArray `gorm:"type:text[]"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	IntegrationID   uuid.UUID      `gorm:"index:idx_notifications_integration"`
	IntegrationType string         `gorm:"index:idx_notifications_integration"`
	EventType       string
	Template        string
	CentralCSURL    string
	CSClusterID     string
	CSClusterName   string
	OwnCSURL        string

	// EmailConfig can only be set in case when IntegrationType == IntegrationEmail
	EmailConfig *EmailConfig `gorm:"type:jsonb"`
	// WebhookConfig can only be set in case when IntegrationType == IntegrationWebhook
	WebhookConfig *WebhookConfig `gorm:"type:jsonb"`
	// SyslogConfig can only be set in case when IntegrationType == IntegrationSyslog
	SyslogConfig *SyslogConfig `gorm:"type:jsonb"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	base := any(&n.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if n.Name != "" {
		if err := n.checkNameUnique(tx, n.Name); err != nil {
			return err
		}
	}

	return nil
}

func (n *Notification) BeforeUpdate(tx *gorm.DB) error {
	base := any(&n.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := n.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (n *Notification) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Notification{Name: name}).
		Not(&Notification{Base: Base{ID: n.ID}}). // this will be needed in future if updates are implemented
		Take(&Notification{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrNotificationNameInUse
}

func (ec *EmailConfig) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, ec)
}

func (ec *EmailConfig) Value() (driver.Value, error) {
	if ec == nil {
		return nil, nil
	}
	return json.Marshal(ec)
}

func (wc *WebhookConfig) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, wc)
}

func (wc *WebhookConfig) Value() (driver.Value, error) {
	if wc == nil {
		return nil, nil
	}
	return json.Marshal(wc)
}

func (sc *SyslogConfig) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, sc)
}

func (sc *SyslogConfig) Value() (driver.Value, error) {
	if sc == nil {
		return nil, nil
	}
	return json.Marshal(sc)
}
