package notifier

import (
	"context"
	"fmt"

	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/notifier/email"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/notifier/syslog"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/notifier/webhook"
)

type Notifier interface {
	Test(context.Context) error
	Notify(ctx context.Context, n *model.Notification, event any) error
}

func FromIntegration(i model.Integration) (Notifier, error) {
	switch conf := i.(type) {
	case *model.Email:
		return &email.Notifier{conf}, nil
	case *model.Webhook:
		return &webhook.Notifier{conf}, nil
	case *model.Syslog:
		return &syslog.Notifier{conf}, nil
	default:
		return nil, fmt.Errorf("invalid integration type given: %T", conf)
	}
}
