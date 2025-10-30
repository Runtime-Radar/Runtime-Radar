package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventTypeRuntimeEvent = "runtime_event"
)

// Event represents a struct that can be extended and used for storing event info.
type Event struct {
	Base
	RegisteredAt time.Time `gorm:"index"`
	Source       string
	Type         string     `gorm:"index"`
	Incident     *Incident  `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	UserID       *uuid.UUID `gorm:"index"`
}
