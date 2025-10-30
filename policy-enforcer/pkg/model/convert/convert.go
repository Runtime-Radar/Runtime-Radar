package convert

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/api"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"gorm.io/gorm"
)

func RulesToProto(rs []*model.Rule) (pbr []*api.Rule) {
	for _, r := range rs {
		pbr = append(pbr, &api.Rule{
			Id:    r.ID.String(),
			Name:  r.Name,
			Rule:  (*api.Rule_RuleJSON)(r.Rule),
			Scope: (*api.Rule_Scope)(r.Scope),
			Type:  RuleTypeToProto(r.Type),
		})
	}

	return
}

func RuleFromProto(pbr *api.Rule) (*model.Rule, error) {
	var id uuid.UUID
	var err error

	if idStr := pbr.GetId(); idStr != "" {
		id, err = uuid.Parse(pbr.GetId())
		if err != nil {
			return nil, fmt.Errorf("can't parse ID: %w", err)
		}
	}

	return &model.Rule{
		model.Base{ID: id},
		pbr.GetName(),
		(*model.RuleJSON)(pbr.GetRule()),
		(*model.Scope)(pbr.GetScope()),
		gorm.DeletedAt{},
		RuleTypeFromProto(pbr.GetType()),
	}, nil
}

func RuleTypeFromProto(pbrt api.Rule_Type) model.RuleType {
	switch pbrt {
	case api.Rule_TYPE_RUNTIME:
		return model.RuleTypeRuntime
	default: // normally should not happen
		panic(fmt.Sprintf("invalid rule type given: %s", pbrt))
	}
}

func RuleTypeToProto(rt model.RuleType) api.Rule_Type {
	switch rt {
	case model.RuleTypeRuntime:
		return api.Rule_TYPE_RUNTIME
	default: // normally should not happen
		panic(fmt.Sprintf("invalid rule type given: %s", rt))
	}
}
