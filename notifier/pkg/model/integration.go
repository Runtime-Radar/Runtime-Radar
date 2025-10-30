package model

import (
	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/lib/security/cipher"
)

const (
	IntegrationEmail   = "email"
	IntegrationWebhook = "webhook"
	IntegrationSyslog  = "syslog"
)

func IntegrationTypeSupported(it string) bool {
	switch it {
	case IntegrationEmail, IntegrationWebhook, IntegrationSyslog:
		return true
	default:
		return false
	}
}

type IntegrationMeta struct {
	CSVersion string
}

// Integration encapsulates configuration of integration with external service (smtp server, webhook, etc.)
type Integration interface {
	GetID() uuid.UUID
	// EncryptSesitive encrypts sensitive data (for example, raw passwords) if needed via given cipher.Crypter
	EncryptSensitive(cipher.Crypter)
	DecryptSensitive(cipher.Crypter)
	// MaskSensitive hides sensitive data before configuration is exposed (for example, returned as API response's body)
	MaskSensitive()
	// SetMeta allows to pass some technical info to notification, such as cs version
	SetMeta(IntegrationMeta)
}
