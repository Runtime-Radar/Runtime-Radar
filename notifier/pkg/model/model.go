package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	// ErrIntegrationNameInUse returned when given integration name is used by another entity in DB.
	ErrIntegrationNameInUse = errors.New("integration name in use")
	// ErrNotificationNameInUse returned when given notification name is used by another entity in DB.
	ErrNotificationNameInUse = errors.New("notification name in use")
	ErrNameFieldsNotUnique   = errors.New("name fields not unique")
)

// Base struct represents basic model to be used by all data structs.
type Base struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid"`

	CreatedAt time.Time `gorm:"index:,sort:desc"`
	UpdatedAt time.Time
}

// BeforeCreate callback sets ID as newly generated UUID. This makes it possible to not install 'uuid-ossp' extension
// with PostgreSQL. If ID was set already, just go ahead and do nothing.
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return nil
}

func getUpdateMapValue[T any](tx *gorm.DB, fieldName string) (T, bool) {
	var val T

	if dest, ok := tx.Statement.Dest.(map[string]any); !ok {
		return val, false
	} else if val, ok = dest[fieldName].(T); !ok {
		return val, false
	}

	return val, true
}
