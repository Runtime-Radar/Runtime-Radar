package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/database/clickhouse"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RuntimeStatsGeneric struct {
	api.UnimplementedRuntimeStatsServer

	StatsRepository clickhouse.StatsRepository
}

func (sg *RuntimeStatsGeneric) CountEvents(ctx context.Context, req *api.RuntimeEventsCounterReq) (*api.Counter, error) {
	from, to, err := timeFromPeriod(req.Period)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time period: %v", err)
	}

	var filter any
	if t := req.GetType(); t != "" {
		if reason, ok := validateRuntimeEventType(t); !ok {
			return nil, status.Error(codes.InvalidArgument, reason)
		}

		filter = gorm.Expr("event_type = ?", t)
	}

	cnt, err := sg.StatsRepository.CountEvents(ctx, from, to, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't count runtime events: %v", err)
	}

	return &api.Counter{Count: int32(cnt)}, nil
}

func validateRuntimeEventType(t string) (reason string, ok bool) {
	switch t {
	case model.RuntimeEventTypeProcessExec,
		model.RuntimeEventTypeProcessExit,
		model.RuntimeEventTypeProcessKprobe,
		model.RuntimeEventTypeProcessTracepoint,
		model.RuntimeEventTypeProcessLoader,
		model.RuntimeEventTypeProcessUprobe:
		return "", true
	default:
		return fmt.Sprintf("invalid runtime event type: %s", t), false
	}
}

func timeFromPeriod(p *api.Period) (from, to time.Time, err error) {
	if p == nil {
		return from, to, errors.New("period is not given")
	}

	pf := p.GetFrom()
	if pf == nil {
		return from, to, errors.New("period.from is not given")
	}
	from = pf.AsTime()

	pt := p.GetTo()
	if pt != nil {
		to = pt.AsTime()
	} else {
		to = time.Now()
	}

	return from, to, nil
}
