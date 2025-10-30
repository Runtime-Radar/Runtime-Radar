package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

const (
	RuleVersion  Version = "1"
	ScopeVersion Version = "1"
)

type (
	RuleJSON api.Rule_RuleJSON
	Scope    api.Rule_Scope
)

// RuleType indicates what type of enforcement Rule will be used for
type RuleType uint8

const (
	// iota + 1 is used instead of iota because gorm considers 0 as zero value and it's impossible to use this value as a criteria in query builder
	RuleTypeImage RuleType = iota + 1
	RuleTypeIAC
	RuleTypeAdmission
	RuleTypeRuntime
	RuleTypeImageMalware
	RuleTypePeriodicScan
)

// ScopeableRuleTypes contain rule types which can be scoped (have a Scope).
var ScopeableRuleTypes = []RuleType{
	RuleTypeImage,
	RuleTypeAdmission,
	RuleTypeRuntime,
	RuleTypeImageMalware,
	RuleTypePeriodicScan,
}

type Rule struct {
	Base
	Name      string         `gorm:"index"`
	Rule      *RuleJSON      `gorm:"type:jsonb"`
	Scope     *Scope         `gorm:"type:jsonb"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Type      RuleType       `gorm:"index"`
}

func (r *Rule) BeforeCreate(tx *gorm.DB) error {
	base := any(&r.Base)
	if b, ok := base.(callbacks.BeforeCreateInterface); ok {
		if err := b.BeforeCreate(tx); err != nil {
			return err
		}
	}

	if r.Name != "" {
		if err := r.checkNameUnique(tx, r.Name); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rule) BeforeUpdate(tx *gorm.DB) error {
	base := any(&r.Base)
	if b, ok := base.(callbacks.BeforeUpdateInterface); ok {
		if err := b.BeforeUpdate(tx); err != nil {
			return err
		}
	}

	name, ok := getUpdateMapValue[string](tx, "Name")
	if ok && name != "" {
		if err := r.checkNameUnique(tx, name); err != nil {
			return err
		}
	}

	return nil
}

func (r *Rule) checkNameUnique(tx *gorm.DB, name string) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&Rule{Name: name}).
		Not(&Rule{Base: Base{ID: r.ID}}).
		Take(&Rule{}).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return fmt.Errorf("can't check if names are in use: %w", err)
	}

	return ErrRuleNameInUse
}

func (r *RuleJSON) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, r)
}

func (r *RuleJSON) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (s *Scope) Scan(src interface{}) error {
	b := src.([]byte)
	return json.Unmarshal(b, s)
}

func (s *Scope) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (rt RuleType) String() string {
	switch rt {
	case RuleTypeImage:
		return "image"
	case RuleTypeIAC:
		return "iac"
	case RuleTypeAdmission:
		return "admission"
	case RuleTypeRuntime:
		return "runtime"
	case RuleTypeImageMalware:
		return "malware"
	case RuleTypePeriodicScan:
		return "periodic_scan"
	default:
		return "unknown"
	}
}
