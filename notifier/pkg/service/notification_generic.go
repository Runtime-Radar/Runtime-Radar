package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/lib/pq"
	history "github.com/runtime-radar/runtime-radar/history-api/pkg/model"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/notifier/api"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/database"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model"
	"github.com/runtime-radar/runtime-radar/notifier/pkg/model/convert"
	pkgtemplate "github.com/runtime-radar/runtime-radar/notifier/pkg/template"
	enforcer_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NotificationGeneric is a basic grpc implementation
type NotificationGeneric struct {
	api.UnimplementedNotificationControllerServer

	NotificationRepository database.NotificationRepository
	IntegrationRepository  database.IntegrationRepository

	RuleController enforcer_api.RuleControllerClient
}

func (ng *NotificationGeneric) Create(ctx context.Context, req *api.Notification) (*api.CreateNotificationResp, error) {
	if reason, ok := ng.validateNotification(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	idStr, integrationIDStr := req.GetId(), req.GetIntegrationId()
	var id uuid.UUID
	var integrationID uuid.UUID
	var err error

	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
		}
	}

	if integrationIDStr != "" {
		integrationID, err = uuid.Parse(integrationIDStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "can't parse integrationID: %v", err)
		}

		if reason, ok := ng.validateIntegration(ctx, req.GetIntegrationType(), integrationID); !ok {
			return nil, status.Error(codes.InvalidArgument, reason)
		}
	}

	n := &model.Notification{
		Base:            model.Base{ID: id},
		Name:            req.GetName(),
		IntegrationType: req.GetIntegrationType(),
		IntegrationID:   integrationID,
		Recipients:      req.GetRecipients(),
		Template:        req.GetTemplate(),
		EventType:       req.GetEventType(),
		CentralCSURL:    req.GetCentralCsUrl(),
		CSClusterID:     req.GetCsClusterId(),
		CSClusterName:   req.GetCsClusterName(),
		OwnCSURL:        req.GetOwnCsUrl(),
	}

	switch c := req.GetConfig().(type) {
	case *api.Notification_Email:
		n.EmailConfig = (*model.EmailConfig)(c.Email)
	case *api.Notification_Webhook:
		n.WebhookConfig = (*model.WebhookConfig)(c.Webhook)
	case *api.Notification_Syslog:
		n.SyslogConfig = (*model.SyslogConfig)(c.Syslog)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid config type given: %T", c)
	}

	if err := ng.NotificationRepository.Add(ctx, n); err != nil {
		if errors.Is(err, model.ErrNotificationNameInUse) {
			return nil, errcommon.StatusWithReason(codes.AlreadyExists, NameMustBeUnique, "name field must be unique").Err()
		}
		return nil, status.Errorf(codes.Internal, "can't add notification: %v", err)
	}

	return &api.CreateNotificationResp{Id: n.ID.String()}, nil
}

func (ng *NotificationGeneric) Read(ctx context.Context, req *api.ReadNotificationReq) (*api.ReadNotificationResp, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	n, err := ng.NotificationRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "notification not found")
		}
		return nil, status.Errorf(codes.Internal, "can't get notification: %v", err)
	}

	return &api.ReadNotificationResp{
		Notification: convert.NotificationToPB(n),
		Deleted:      n.DeletedAt.Valid,
	}, nil
}

func (ng *NotificationGeneric) Update(ctx context.Context, req *api.Notification) (*emptypb.Empty, error) {
	if reason, ok := ng.validateNotification(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	var integrationID uuid.UUID
	if req.GetIntegrationId() != "" {
		integrationID, err = uuid.Parse(req.GetIntegrationId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "can't parse integrationID: %v", err)
		}

		if reason, ok := ng.validateIntegration(ctx, req.GetIntegrationType(), integrationID); !ok {
			return nil, status.Error(codes.InvalidArgument, reason)
		}
	}

	m := map[string]any{
		"Name":            req.GetName(),
		"IntegrationType": req.GetIntegrationType(),
		"IntegrationID":   integrationID,
		"Recipients":      pq.StringArray(req.GetRecipients()),
		"Template":        req.GetTemplate(),
		"EventType":       req.GetEventType(),
		"CentralCSURL":    req.GetCentralCsUrl(),
		"CSClusterID":     req.GetCsClusterId(),
		"CSClusterName":   req.GetCsClusterName(),
		"OwnCSURL":        req.GetOwnCsUrl(),
	}

	switch c := req.GetConfig().(type) {
	case *api.Notification_Email:
		m["EmailConfig"] = (*model.EmailConfig)(c.Email)
	case *api.Notification_Webhook:
		m["WebhookConfig"] = (*model.WebhookConfig)(c.Webhook)
	case *api.Notification_Syslog:
		m["SyslogConfig"] = (*model.SyslogConfig)(c.Syslog)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid config type given: %T", c)
	}

	if err := ng.NotificationRepository.UpdateWithMap(ctx, id, m); err != nil {
		if errors.Is(err, model.ErrNotificationNameInUse) {
			return nil, errcommon.StatusWithReason(codes.AlreadyExists, NameMustBeUnique, "name field must be unique").Err()
		}
		return nil, status.Errorf(codes.Internal, "can't update notification: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (ng *NotificationGeneric) Delete(ctx context.Context, req *api.DeleteNotificationReq) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "notification ID is empty")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	rcResp, err := ng.RuleController.NotifyTargetsInUse(ctx, &enforcer_api.NotifyTargetsInUseReq{
		Targets: []string{id.String()},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't check if notification is in use: %v", err)
	}

	if rcResp.InUse {
		return nil, errcommon.StatusWithReason(codes.FailedPrecondition, NotificationInUse, "notification is in use by rule").Err()
	}

	if err := ng.NotificationRepository.DeleteByID(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "can't delete notification: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (ng *NotificationGeneric) List(ctx context.Context, req *api.ListNotificationReq) (*api.ListNotificationResp, error) {
	order := req.GetOrder()
	if order == "" {
		order = defaultOrder
	}

	var filter any

	expr, err := buildNotificationFilter(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't build notification filter: %v", err)
	}

	// when passing clause.Expr with empty sql gorm generates incorrect query
	if expr.SQL != "" {
		filter = expr
	}

	ns, err := ng.NotificationRepository.GetAll(ctx, filter, order)
	if err != nil {
		if errors.Is(err, database.ErrInvalidOrder) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "can't get notifications: %v", err)
	}

	return &api.ListNotificationResp{Notifications: convert.NotificationsToPB(ns)}, nil
}

func (ng *NotificationGeneric) DefaultTemplate(_ context.Context, req *api.DefaultTemplateReq) (*api.DefaultTemplateResp, error) {
	if reason, ok := validateIntegrationType(req.IntegrationType); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	if !validateEventType(req.EventType) {
		return nil, status.Errorf(codes.InvalidArgument, "event type is invalid: %s", req.EventType)
	}

	var (
		filePath string
		ok       bool
	)
	switch req.IntegrationType {
	case model.IntegrationEmail:
		filePath, ok = pkgtemplate.HTMLFilePaths[req.EventType]
	case model.IntegrationWebhook, model.IntegrationSyslog:
		filePath, ok = pkgtemplate.TextFilePaths[req.EventType]
	default:
		return nil, status.Errorf(codes.Internal, "unexpected integration type: %s", req.IntegrationType)
	}

	if !ok {
		return nil, status.Errorf(codes.NotFound, "default template for '%s' event type not found", req.EventType)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't read default template: %v", err)
	}

	return &api.DefaultTemplateResp{
		Template: string(bytes),
	}, nil
}

func buildNotificationFilter(req *api.ListNotificationReq) (clause.Expr, error) {
	exprs := []string{}
	args := []any{}

	if req.EventType != "" {
		exprs = append(exprs, "event_type = ?")
		args = append(args, req.EventType)
	}

	if req.IntegrationId != "" {
		integrationID, err := uuid.Parse(req.IntegrationId)
		if err != nil {
			return clause.Expr{}, err
		}

		exprs = append(exprs, "integration_id = ?")
		args = append(args, integrationID)
	}

	sql := strings.Join(exprs, " and ")

	return gorm.Expr(sql, args...), nil
}

func (ng *NotificationGeneric) validateNotification(n *api.Notification) (string, bool) {
	if n.GetName() == "" {
		return "notification name is empty", false
	}
	if n.GetIntegrationId() == "" {
		return "notification integration ID is empty", false
	}
	if et := n.GetEventType(); et == "" {
		return "event type is empty", false
	} else if !validateEventType(et) {
		return fmt.Sprintf("event type is invalid: %s", et), false
	}

	switch it := n.GetIntegrationType(); it {
	case model.IntegrationEmail:
		if len(n.GetRecipients()) == 0 {
			return "notification recipients is empty", false
		}

		if c := n.GetConfig(); c != nil {
			emailConf, ok := c.(*api.Notification_Email)
			if !ok {
				return fmt.Sprintf("integration type %s mismatches config type %T", it, c), false
			}

			if subjTpl := emailConf.Email.GetSubjectTemplate(); subjTpl != "" {
				if _, err := pkgtemplate.NewText("", subjTpl); err != nil {
					return fmt.Sprintf("invalid subject template: %v", err), false
				}
			}
		}
	case model.IntegrationWebhook:
		if c := n.GetConfig(); c != nil {
			if _, ok := c.(*api.Notification_Webhook); !ok {
				return fmt.Sprintf("integration type %s mismatches config type %T", it, c), false
			}
		}
	case model.IntegrationSyslog:
		if c := n.GetConfig(); c != nil {
			if _, ok := c.(*api.Notification_Syslog); !ok {
				return fmt.Sprintf("integration type %s mismatches config type %T", it, c), false
			}
		}
	default:
		return fmt.Sprintf("unsupported integration type given: %s", n.GetIntegrationType()), false
	}

	if n.GetTemplate() != "" {
		tpl, err := pkgtemplate.NewText("", n.GetTemplate())
		if err != nil {
			return fmt.Sprintf("invalid template %v", err), false
		}

		// in case of webhook/syslog integration we render template with fake data and check that result is valid json
		if n.GetIntegrationType() == model.IntegrationWebhook || n.GetIntegrationType() == model.IntegrationSyslog {
			if reason, ok := validateTemplate(n.GetEventType(), tpl); !ok {
				return reason, false
			}
		}
	}

	return "", true
}

func (ng *NotificationGeneric) validateIntegration(ctx context.Context, it string, integrationID uuid.UUID) (string, bool) {
	if _, err := ng.IntegrationRepository.GetByTypeAndID(ctx, it, integrationID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "can't find integration", false
		}
		return "can't get integration", false
	}
	return "", true
}

func validateEventType(et string) bool {
	return et == history.EventTypeRuntimeEvent
}

func validateTemplate(eventType string, tpl *template.Template) (string, bool) {
	data := map[string]any{
		"notificationName": "notification name",
		"event":            fakeScanEvent(eventType),
	}

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return fmt.Sprintf("can't execute template: %v", err), false
	}

	if !json.Valid(buf.Bytes()) {
		return "template execution's result is not valid json", false
	}

	return "", true
}

func fakeScanEvent(eventType string) any {
	switch eventType {
	case history.EventTypeRuntimeEvent:
		pid := uint32(2)
		uid := uint32(0)
		parentPID := uint32(1)
		parentUID := uint32(0)
		setuid := uint32(0)

		return &api.RuntimeEvent{
			Event: &api.RuntimeEvent_Event{
				Threats: []*api.RuntimeEvent_Event_Threat{
					{
						DetectorId:          "POSI_AC_001",
						DetectorName:        "My Detector",
						DetectorVersion:     1,
						DetectorDescription: "Very dangerous security issue",
						Severity:            "high",
					},
				},
				EventType:           "PROCESS_EXEC",
				PodNamespace:        "default",
				PodName:             "deathstar-8464cdd4d9-7slzs",
				ContainerName:       "deathstar",
				ContainerImage:      "docker.io/cilium/starwars:latest",
				ContainerId:         "deathstar-id",
				FunctionName:        "func",
				ProcessBinary:       "/usr/bin/nc",
				ProcessArguments:    "-vz postgres 5432",
				ProcessPid:          &pid,
				ProcessUid:          &uid,
				ProcessCapEffective: []string{"CAP_NET_RAW"},
				ProcessCapPermitted: []string{"CAP_NET_RAW"},
				ProcessSetuid:       &setuid,
				ProcessSetgid:       nil,
				ParentPid:           &parentPID,
				ParentUid:           &parentUID,
				ParentBinary:        "/bin/sh",
				ParentArguments:     "-c /usr/bin/nc",
				NodeName:            "node001",
			},
			RegisteredAt: timestamppb.Now(),
			Severity:     "high",
			Block:        true,
			RuleName:     "rule 1",
			EventId:      uuid.NewString(),
		}
	default: // normally should not happen
		panic(fmt.Sprintf("invalid event type: %s", eventType))
	}
}
