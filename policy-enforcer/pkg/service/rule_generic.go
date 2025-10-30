package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/runtime-radar/runtime-radar/lib/errcommon"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/cache"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/database"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model/convert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var (
	// see https://github.com/distribution/distribution/blob/v2.7.1/reference/reference.go#L4
	imagePatternRegex = regexp.MustCompile(`^[a-zA-Z0-9._:+@*\-?/]+$`)

	// see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
	namespacePatternRegex = regexp.MustCompile("^[a-z0-9-*?]+$")
	podPatternRegex       = regexp.MustCompile("^[a-z0-9-*?]+$")
	containerPatternRegex = regexp.MustCompile("^[a-z0-9-*?]+$")

	// see https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
	nodePatternRegex    = regexp.MustCompile("^[a-z0-9-.*?]+$")
	clusterPatternRegex = regexp.MustCompile("^[a-z0-9-.*?]+$")

	allowedSeverities = map[string]bool{
		model.NoneSeverity.String():     true,
		model.LowSeverity.String():      true,
		model.MediumSeverity.String():   true,
		model.HighSeverity.String():     true,
		model.CriticalSeverity.String(): true,
	}
)

// RuleGeneric is basic grpc service implementation.
type RuleGeneric struct {
	api.UnimplementedRuleControllerServer

	RuleRepository database.RuleRepository
	RuleMatcher    cache.RuleMatcher
}

func (rg *RuleGeneric) Create(ctx context.Context, req *api.Rule) (*api.CreateRuleResp, error) {
	if reason, ok := rg.validateRule(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	r, err := convert.RuleFromProto(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't convert rule: %v", err)
	}

	if err := rg.RuleRepository.Add(ctx, r); err != nil {
		if errors.Is(err, model.ErrRuleNameInUse) {
			return nil, errcommon.StatusWithReason(codes.AlreadyExists, NameMustBeUnique, "name field must be unique").Err()
		}
		return nil, status.Errorf(codes.Internal, "can't add rule: %v", err)
	}

	if err := rg.RuleMatcher.Invalidate(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "can't invalidate scope cache: %v", err)
	}

	resp := &api.CreateRuleResp{
		Id: r.ID.String(),
	}

	return resp, nil
}

func (rg *RuleGeneric) Read(ctx context.Context, req *api.ReadRuleReq) (*api.ReadRuleResp, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	r, err := rg.RuleRepository.GetByID(ctx, id, true) // preload is on
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "rule not found")
		}
		return nil, status.Errorf(codes.Internal, "can't read rule: %v", err)
	}

	resp := &api.ReadRuleResp{
		Rule: &api.Rule{
			Id:    r.ID.String(),
			Name:  r.Name,
			Rule:  (*api.Rule_RuleJSON)(r.Rule),
			Scope: (*api.Rule_Scope)(r.Scope),
			Type:  convert.RuleTypeToProto(r.Type),
		},
		Deleted: r.DeletedAt.Valid,
	}

	return resp, nil
}

func (rg *RuleGeneric) Update(ctx context.Context, req *api.Rule) (*emptypb.Empty, error) {
	if reason, ok := rg.validateRule(req); !ok {
		return nil, status.Error(codes.InvalidArgument, reason)
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	m := map[string]any{
		"Name":  req.GetName(),
		"Rule":  (*model.RuleJSON)(req.GetRule()),
		"Scope": (*model.Scope)(req.GetScope()),
	}

	if err := rg.RuleRepository.UpdateWithMap(ctx, id, m); err != nil {
		if errors.Is(err, model.ErrRuleNameInUse) {
			return nil, errcommon.StatusWithReason(codes.AlreadyExists, NameMustBeUnique, "name field must be unique").Err()
		}
		return nil, status.Errorf(codes.Internal, "can't update rule: %v", err)
	}

	if err := rg.RuleMatcher.Invalidate(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "can't invalidate scope cache: %v", err)
	}

	resp := &emptypb.Empty{}

	return resp, nil
}

func (rg *RuleGeneric) Delete(ctx context.Context, req *api.DeleteRuleReq) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse ID: %v", err)
	}

	if err := rg.RuleRepository.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "can't delete rule: %v", err)
	}

	if err := rg.RuleMatcher.Invalidate(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "can't invalidate scope cache: %v", err)
	}

	resp := &emptypb.Empty{}

	return resp, nil
}

func (rg *RuleGeneric) ListPage(ctx context.Context, req *api.ListRulePageReq) (*api.ListRulePageResp, error) {
	total, err := rg.RuleRepository.GetCount(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get rule count: %v", err)
	}

	order, pageSize := req.GetOrder(), req.GetPageSize()
	if order == "" {
		order = defaultOrder
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	rs, err := rg.RuleRepository.GetPage(ctx, nil, order, int(pageSize), int(req.GetPageNum()), true) // preload is on
	if err != nil {
		if errors.Is(err, database.ErrInvalidOrder) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "can't get rule page: %v", err)
	}

	resp := &api.ListRulePageResp{
		Total: uint32(total),
		Rules: convert.RulesToProto(rs),
	}

	return resp, nil
}

func (rg *RuleGeneric) NotifyTargetsInUse(ctx context.Context, req *api.NotifyTargetsInUseReq) (*api.NotifyTargetsInUseResp, error) {
	if len(req.Targets) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no targets given")
	}

	filter := gorm.Expr(`rule->'notify'->'targets' ? ?`, gorm.Expr("?|"), pq.StringArray(req.Targets))

	count, err := rg.RuleRepository.GetCount(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can't get rule count: %v", err)
	}

	return &api.NotifyTargetsInUseResp{InUse: count > 0}, nil
}

func (rg *RuleGeneric) validateRule(req *api.Rule) (reason string, ok bool) {
	if req.GetName() == "" {
		return "empty or missing name", false
	} else if req.Rule == nil {
		return "no rule", false
	}

	if req.Rule.GetVersion() == "" {
		return "empty or missing rule version", false
	} else if ver := req.Rule.GetVersion(); ver != string(model.RuleVersion) {
		return fmt.Sprintf("rule version mismatch: expected %s, got %s", model.RuleVersion, ver), false
	}

	if req.Rule.Block == nil && req.Rule.Notify == nil {
		return "at least one block or notify required", false
	}

	if req.Rule.Block != nil {
		if req.Rule.Block.GetSeverity() == "" {
			return "at least one: severity or verdict fields are required in block", false
		}

		if req.Rule.Block.GetSeverity() != "" && !allowedSeverities[req.Rule.Block.GetSeverity()] {
			return "severity must be one of [none|low|medium|high|critical] in block", false
		}
	}

	if req.Rule.Notify != nil {
		if req.Rule.Notify.GetSeverity() == "" {
			return "at least one: severity or verdict fields are required in notify", false
		}

		if req.Rule.Notify.GetSeverity() != "" && !allowedSeverities[req.Rule.Notify.GetSeverity()] {
			return "severity must be one of [none|low|medium|high|critical] in notify", false
		}
	}

	for _, b := range req.Rule.Whitelist.GetBinaries() {
		if _, err := glob.Compile(b); err != nil {
			return fmt.Sprintf("can't compile whitelist binary name pattern: %v", err), false
		}
	}

	s := req.GetScope()
	t := req.GetType()

	if s == nil {
		return fmt.Sprintf("scope is empty for rule type requiring scope: %s", t), false
	} else if s != nil {
		if reason, ok := rg.validateScope(s); !ok {
			return fmt.Sprintf("scope is invalid: %s", reason), false
		}
	}

	return "", true
}

func (rg *RuleGeneric) validateScope(s *api.Rule_Scope) (reason string, ok bool) {
	if v := s.GetVersion(); v == "" {
		return "empty or missing version", false
	} else if v != string(model.ScopeVersion) {
		return fmt.Sprintf("version mismatch: expected %s, got %s", model.ScopeVersion, v), false
	}

	for _, img := range s.GetImageNames() {
		if !imagePatternRegex.MatchString(img) {
			return fmt.Sprintf("invalid image name: %s", img), false
		}
	}

	for _, ns := range s.GetNamespaces() {
		if !namespacePatternRegex.MatchString(ns) {
			return fmt.Sprintf("invalid namespace name: %s", ns), false
		}
	}

	for _, c := range s.GetClusters() {
		if !clusterPatternRegex.MatchString(c) {
			return fmt.Sprintf("invalid cluster name: %s", c), false
		}
	}

	for _, p := range s.GetPods() {
		if !podPatternRegex.MatchString(p) {
			return fmt.Sprintf("invalid pod name: %s", p), false
		}
	}

	for _, c := range s.GetContainers() {
		if !containerPatternRegex.MatchString(c) {
			return fmt.Sprintf("invalid container name: %s", c), false
		}
	}

	for _, n := range s.GetNodes() {
		if !nodePatternRegex.MatchString(n) {
			return fmt.Sprintf("invalid node name: %s", n), false
		}
	}

	return "", true
}
