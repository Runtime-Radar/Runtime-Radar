package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/cluster-manager/api"
	"github.com/runtime-radar/runtime-radar/lib/security/cipher"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

const (
	ClusterConfigVersion Version = "1"
)

type ClusterStatus uint8

const (
	ClusterStatusUnregistered = iota
	ClusterStatusRegistered
)

type ClusterConfig api.Cluster_Config

type Cluster struct {
	Base
	Name         string         `gorm:"index"`
	Token        uuid.UUID      `gorm:"uniqueIndex"`
	Status       ClusterStatus  `gorm:"index"`
	Config       *ClusterConfig `gorm:"type:jsonb"`
	RegisteredAt *time.Time     `gorm:"index"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (c *Cluster) BeforeCreate(tx *gorm.DB) error {
	base := any(&c.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if c.Name != "" {
		if err := c.checkNameUnique(tx, c.Name); err != nil {
			return err
		}
	}

	if c.Token == uuid.Nil {
		c.Token = uuid.New()
	}

	return nil
}

func (c *Cluster) BeforeUpdate(tx *gorm.DB) error {
	base := any(&c.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := c.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Cluster{Name: name}).
		Not(&Cluster{Base: Base{ID: c.ID}}).
		Take(&Cluster{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrClusterNameInUse
}

func (cs ClusterStatus) String() string {
	switch cs {
	case ClusterStatusUnregistered:
		return "unregistered"
	case ClusterStatusRegistered:
		return "registered"
	default:
		return "unknown"
	}
}

func (cc *ClusterConfig) Scan(src any) error {
	b := src.([]byte)
	return json.Unmarshal(b, cc)
}

func (cc *ClusterConfig) Value() (driver.Value, error) {
	return json.Marshal(cc)
}

func (cc *ClusterConfig) EncryptSensitive(c cipher.Crypter) {
	if pg := cc.Postgres; pg != nil && pg.Password != "" {
		cc.Postgres.Password = c.EncryptStringAsHex(pg.Password)
	}

	if ch := cc.Clickhouse; ch != nil && ch.Password != "" {
		cc.Clickhouse.Password = c.EncryptStringAsHex(ch.Password)
	}

	if r := cc.Redis; r != nil && r.Password != "" {
		cc.Redis.Password = c.EncryptStringAsHex(r.Password)
	}

	if r := cc.Rabbit; r != nil && r.Password != "" {
		cc.Rabbit.Password = c.EncryptStringAsHex(r.Password)
	}

	if r := cc.Registry; r != nil && r.Password != "" {
		cc.Registry.Password = c.EncryptStringAsHex(r.Password)
	}
}

func (cc *ClusterConfig) DecryptSensitive(c cipher.Crypter) {
	if pg := cc.Postgres; pg != nil && pg.Password != "" {
		cc.Postgres.Password = c.DecryptHexAsString(pg.Password)
	}

	if ch := cc.Clickhouse; ch != nil && ch.Password != "" {
		cc.Clickhouse.Password = c.DecryptHexAsString(ch.Password)
	}

	if r := cc.Redis; r != nil && r.Password != "" {
		cc.Redis.Password = c.DecryptHexAsString(r.Password)
	}

	if r := cc.Rabbit; r != nil && r.Password != "" {
		cc.Rabbit.Password = c.DecryptHexAsString(r.Password)
	}

	if r := cc.Registry; r != nil && r.Password != "" {
		cc.Registry.Password = c.DecryptHexAsString(r.Password)
	}
}

func (cc *ClusterConfig) MaskSensitive() {
	const mask = "********"

	if pg := cc.Postgres; pg != nil && pg.Password != "" {
		cc.Postgres.Password = mask
	}

	if ch := cc.Clickhouse; ch != nil && ch.Password != "" {
		cc.Clickhouse.Password = mask
	}

	if r := cc.Redis; r != nil && r.Password != "" {
		cc.Redis.Password = mask
	}

	if r := cc.Rabbit; r != nil && r.Password != "" {
		cc.Rabbit.Password = mask
	}

	if r := cc.Registry; r != nil && r.Password != "" {
		cc.Registry.Password = mask
	}
}
