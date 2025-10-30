// nolint: goconst
package postgres

import (
	"context"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/lib/security"
	enf_model "github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"gorm.io/gorm"
)

func populate(db *gorm.DB, count int) error {
	if count == 0 {
		return nil
	}

	eventDB := &EventDatabase{db}

	es := []*model.Event{} // some events can be added manually
	es = addEvents(es, count)

	return eventDB.Add(context.Background(), es...)
}

// nolint:gosec
func addEvents(es []*model.Event, count int) []*model.Event {
	for i := count; i >= 1; i-- {
		createdAt := time.Now().Add(time.Duration(-i) * time.Hour)
		registeredAt := createdAt

		severity := enf_model.Severity(rand.Intn(5))

		e := &model.Event{
			Base:         model.Base{CreatedAt: createdAt},
			RegisteredAt: registeredAt,
			Source:       "cli",
			Type:         "image_scan",
			Incident: &model.Incident{
				Base:     model.Base{CreatedAt: createdAt},
				BlockBy:  model.BlockBy{"74875022-32b9-4089-a665-63c230ed4d63", "74875022-32b9-4089-a665-63c230ed4d68"},
				NotifyBy: model.NotifyBy{"74875022-32b9-4089-a665-63c230ed4d63", "74875022-32b9-4089-a665-63c230ed4d68"},
				Severity: severity,
			},
		}

		es = append(es, e)
	}

	return es
}

func randDigest(prefix string, size int) string {
	return prefix + ":" + hex.EncodeToString(security.Rand(size))
}
