package service

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	processor_api "github.com/runtime-radar/runtime-radar/event-processor/api"
	"github.com/runtime-radar/runtime-radar/history-api/api"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/database/clickhouse"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/history-api/pkg/model/convert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type RuntimeHistoryGeneric struct {
	api.UnimplementedRuntimeHistoryServer

	RuntimeEventRepository clickhouse.RuntimeEventRepository
}

func (rhg *RuntimeHistoryGeneric) Read(ctx context.Context, req *api.ReadRuntimeEventReq) (*processor_api.RuntimeEvent, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	e, err := rhg.RuntimeEventRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "runtime event not found")
		}
		return nil, status.Errorf(codes.Internal, "can't read runtime event: %v", err)
	}

	resp, err := convert.RuntimeEventToProto(e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't convert runtime event model to proto: %v", err)
	}

	return resp, nil
}

func (rhg *RuntimeHistoryGeneric) ListRuntimeEventSlice(ctx context.Context, req *api.ListRuntimeEventSliceReq) (*api.ListRuntimeEventSliceResp, error) {
	if reason, ok := rhg.validateListRuntimeEventSliceReq(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	sliceSize := int(req.GetSliceSize())
	if sliceSize == 0 {
		sliceSize = defaultSliceSize
	}

	var events []*model.RuntimeEvent
	var err error

	if req.GetDirection() == directionLeft {
		events, err = rhg.RuntimeEventRepository.GetLeftSlice(ctx, req.GetCursor().AsTime(), nil, sliceSize)
	} else {
		events, err = rhg.RuntimeEventRepository.GetRightSlice(ctx, req.GetCursor().AsTime(), nil, sliceSize)
	}

	if err != nil {
		return nil, fmt.Errorf("can't get events from db: %w", err)
	}

	runtimeEvents, err := convert.RuntimeEventsToProto(events)
	if err != nil {
		return nil, err
	}

	var leftCursor, rightCursor *timestamppb.Timestamp
	if len(events) != 0 {
		leftCursor = timestamppb.New(events[0].RegisteredAt)
		rightCursor = timestamppb.New(events[len(events)-1].RegisteredAt)
	}

	resp := &api.ListRuntimeEventSliceResp{
		LeftCursor:    leftCursor,
		RightCursor:   rightCursor,
		RuntimeEvents: runtimeEvents,
	}

	return resp, nil
}

func (rhg *RuntimeHistoryGeneric) FilterRuntimeEventSlice(ctx context.Context, req *api.FilterRuntimeEventSliceReq) (*api.ListRuntimeEventSliceResp, error) {
	if reason, ok := rhg.validateFilterRuntimeEventsReq(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	sliceSize := int(req.GetSliceSize())
	if sliceSize == 0 {
		sliceSize = defaultSliceSize
	}

	filter, err := makeRuntimeEventFilter(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't make runtime filter: %v", err)
	}

	var events []*model.RuntimeEvent

	if req.GetDirection() == directionLeft {
		events, err = rhg.RuntimeEventRepository.GetLeftSlice(ctx, req.GetCursor().AsTime(), filter, sliceSize)
	} else {
		events, err = rhg.RuntimeEventRepository.GetRightSlice(ctx, req.GetCursor().AsTime(), filter, sliceSize)
	}

	if err != nil {
		return nil, fmt.Errorf("can't get events from db: %w", err)
	}

	runtimeEvents, err := convert.RuntimeEventsToProto(events)
	if err != nil {
		return nil, err
	}

	var leftCursor, rightCursor *timestamppb.Timestamp
	if len(events) != 0 {
		leftCursor = timestamppb.New(events[0].RegisteredAt)
		rightCursor = timestamppb.New(events[len(events)-1].RegisteredAt)
	}

	resp := &api.ListRuntimeEventSliceResp{
		LeftCursor:    leftCursor,
		RightCursor:   rightCursor,
		RuntimeEvents: runtimeEvents,
	}

	return resp, nil
}

// nolint:goconst
func (rhg *RuntimeHistoryGeneric) validateListRuntimeEventSliceReq(req *api.ListRuntimeEventSliceReq) (string, bool) {
	if req.GetDirection() != directionLeft && req.GetDirection() != directionRight {
		return fmt.Sprintf("unsupported direction: %s", req.GetDirection()), false
	}

	cursor := req.GetCursor()
	if cursor == nil {
		return "missing cursor", false
	}

	if err := cursor.CheckValid(); err != nil {
		return fmt.Sprintf("cursor is invalid: %s", err.Error()), false
	}

	return "", true
}

func (rhg *RuntimeHistoryGeneric) validateFilterRuntimeEventsReq(req *api.FilterRuntimeEventSliceReq) (reason string, ok bool) {
	if req.GetDirection() != directionLeft && req.GetDirection() != directionRight {
		return fmt.Sprintf("unsupported direction: %s", req.GetDirection()), false
	}

	cursor := req.GetCursor()
	if cursor == nil {
		return "missing cursor", false
	}

	if err := cursor.CheckValid(); err != nil {
		return fmt.Sprintf("cursor is invalid: %s", err.Error()), false
	}

	if reason, ok := rhg.validateRuntimeEventFilter(req.GetFilter()); !ok {
		return fmt.Sprintf("filter is invalid: %s", reason), ok
	}

	return "", true
}

func (rhg *RuntimeHistoryGeneric) validateRuntimeEventFilter(rf *api.RuntimeFilter) (reason string, ok bool) {
	if len(rf.GetEventType()) == 0 &&
		len(rf.GetKprobeFunctionName()) == 0 &&
		len(rf.GetProcessPodNamespace()) == 0 &&
		len(rf.GetProcessPodName()) == 0 &&
		len(rf.GetNodeName()) == 0 &&
		len(rf.GetProcessPodContainerName()) == 0 &&
		len(rf.GetProcessPodContainerImageName()) == 0 &&
		len(rf.GetProcessBinary()) == 0 &&
		len(rf.GetProcessArguments()) == 0 &&
		rf.GetPeriod().GetFrom() == nil &&
		rf.GetPeriod().GetTo() == nil &&
		rf.HasThreats == nil &&
		rf.GetProcessExecId() == "" &&
		rf.GetProcessParentExecId() == "" &&
		len(rf.GetThreatsDetectors()) == 0 &&
		len(rf.GetRules()) == 0 &&
		rf.HasIncident == nil {
		return "at least one filter should be set", false
	}

	for _, t := range rf.GetEventType() {
		if reason, ok := rhg.validateRuntimeEventType(t); !ok {
			return reason, ok
		}
	}

	if len(rf.KprobeFunctionName) > 0 {
		if et := rf.GetEventType(); len(et) > 0 && !slices.Contains(et, model.RuntimeEventTypeProcessKprobe) {
			return fmt.Sprintf("kprobe function name is passed when type doesn't contain %s", model.RuntimeEventTypeProcessKprobe), false
		}
	}

	for i, r := range rf.GetRules() {
		if _, err := uuid.Parse(r); err != nil {
			return fmt.Sprintf("can't parse rules[%d]: %v", i, err), false
		}
	}

	return "", true
}

func (rhg *RuntimeHistoryGeneric) validateRuntimeEventType(t string) (reason string, ok bool) {
	switch t {
	case model.RuntimeEventTypeUndef,
		model.RuntimeEventTypeProcessExec,
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
