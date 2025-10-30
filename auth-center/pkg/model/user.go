package model

import (
	"time"

	"github.com/google/uuid"
)

const AdminInitIDStr = "00000000-0000-0000-0000-000000000001"

var AdminInitID = uuid.MustParse(AdminInitIDStr)

type AuthType string

var (
	AuthTypeInternal AuthType = "internal"
	AuthTypeLDAP     AuthType = "ldap"
)

func (a AuthType) String() string {
	return string(a)
}

func ValidateAuthType(a string) (AuthType, bool) {
	if a == string(AuthTypeInternal) || a == string(AuthTypeLDAP) {
		return AuthType(a), true
	}

	return AuthTypeInternal, true
}

type User struct {
	Base
	Username              string `gorm:"unique;not null"`
	Email                 string
	AuthType              AuthType  `gorm:"not null; default='internal';"`
	HashedPassword        string    `gorm:"not null;"`
	RoleID                uuid.UUID `gorm:"foreignKey:RoleID; not null;"`
	Role                  Role      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MappingRoleID         *uuid.UUID
	LastPasswordChangedAt time.Time `gorm:"not null"`
}
