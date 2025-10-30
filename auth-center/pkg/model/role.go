package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/runtime-radar/runtime-radar/lib/security/jwt"
)

const AdminRoleIDStr = "00000000-0000-0000-0000-000000000001"

var AdminRoleID = uuid.MustParse(AdminRoleIDStr)

type (
	Permissions jwt.RolePermissions
)

func (p *Permissions) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, p)
}

func (p *Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type RolePermission struct {
	Actions     pq.StringArray `gorm:"type:text[]"` // Утончить по такому использованию
	Description string
}

type Role struct {
	ID              uuid.UUID   `gorm:"primaryKey; type:uuid" json:"id"`
	RoleName        string      `json:"role_name"`
	RolePermissions Permissions `gorm:"type:jsonb" json:"role_permissions"`
	Description     string      `json:"description"`
}
