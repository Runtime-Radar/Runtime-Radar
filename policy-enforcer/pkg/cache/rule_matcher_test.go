package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/database"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
)

type ruleRepositoryMock struct {
	database.RuleRepository
	rules []*model.Rule
}

func (rm *ruleRepositoryMock) GetAll(_ context.Context, _, _ interface{}, _ bool) ([]*model.Rule, error) {
	return rm.rules, nil
}

func TestMatchRules(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		Name          string
		Rules         []*model.Rule
		RuleType      model.RuleType
		Opts          []MatchOption
		Matches       int
		ExpectedError error
	}{
		{
			"All",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("any"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Myimage",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Myimage",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "myimage:1.0.8"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Mygroup",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Mygroup with question mark",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup with question mark",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/myimage:1.0.?"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Mygroup with character range",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup with character range",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/myimage:1.0.[8-9]"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Not my group",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Not my group",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("notmygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			0,
			nil,
		},
		{
			"Mygroup and mynamespace",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup and mynamespace",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"mynamespace-*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("mynamespace-dev"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Mygroup and mynamespace and mycluster",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup and mynamespace and mycluster",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"mynamespace-*"},
						Clusters:   []string{"mycluster-*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("mynamespace-dev"),
				WithCluster("mycluster-test"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Mygroup and mynamespace and not mycluster",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "Mygroup and mynamespace and not mycluster",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"mynamespace-*"},
						Clusters:   []string{"mycluster-*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("mynamespace-dev"),
				WithCluster("notmycluster-test"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			0,
			nil,
		},
		{
			"Two matches with myimage",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")},
					Name: "mygroup",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"qwerty", "mygroup/*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			2,
			nil,
		},
		{
			"My registry",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"registry-*.example.com"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("registry-docker.example.com"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			1,
			nil,
		},
		{
			"Not my registry",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"registry-*.example.com"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("mygroup/myimage:1.0.8"),
				WithRegistry("some-registry.com"),
				WithNamespace("any"),
				WithCluster("any"),
				WithPod("any"),
				WithContainer("any"),
				WithNode("any"),
			},
			0,
			nil,
		},
		{
			"Without cluster, pod and node",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{},
						Registries: []string{"*"},
						Pods:       []string{},
						Containers: []string{"*"},
						Nodes:      []string{},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithImageName("myimage:1.0.8"),
				WithRegistry("any"),
				WithNamespace("any"),
				WithContainer("any"),
			},
			1,
			nil,
		},
		{
			"Without container, image and registry",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{},
						Pods:       []string{"*"},
						Containers: []string{},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithNamespace("any"),
				WithCluster("any"),
			},
			1,
			nil,
		},
		{
			"With image and registry in rule's scope but without them in options",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"my-image-?"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*-my-registry.com"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{
				WithNamespace("any"),
				WithCluster("any"),
			},
			1,
			nil,
		},
		{
			"Matching by namespace but mismatching by type",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeImage,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")},
					Name: "mygroup",
					Type: model.RuleTypeRuntime,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"my-namespace*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeAdmission,
			[]MatchOption{
				WithNamespace("my-namespace-1"),
			},
			0,
			nil,
		},
		{
			"All rules matching by namespace but only one matching by type",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeRuntime,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002")},
					Name: "mygroup",
					Type: model.RuleTypeAdmission,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"my-namespace*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeAdmission,
			[]MatchOption{
				WithNamespace("my-namespace-2"),
			},
			1,
			nil,
		},
		{
			"Not scopeable rule type",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeRuntime,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeIAC, // at the moment model.RuleTypeIAC is the only type of rule that is not scopeable. See model.ScopeableRuleTypes.
			[]MatchOption{
				WithNamespace("irrelevant-value"),
			},
			0,
			ErrNotScopeable,
		},
		{
			"No options passed to matcher",
			[]*model.Rule{
				{
					Base: model.Base{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001")},
					Name: "All",
					Type: model.RuleTypeRuntime,
					Scope: &model.Scope{
						Version:    string(model.ScopeVersion),
						ImageNames: []string{"*"},
						Namespaces: []string{"*"},
						Clusters:   []string{"*"},
						Registries: []string{"*"},
						Pods:       []string{"*"},
						Containers: []string{"*"},
						Nodes:      []string{"*"},
					},
				},
			},
			model.RuleTypeImage,
			[]MatchOption{},
			0,
			ErrNoOptions,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			local := NewLocal()
			repoMock := &ruleRepositoryMock{
				rules: tc.Rules,
			}

			rc := &RuleCache{
				local,
				repoMock,
			}

			rs, err := rc.MatchRules(context.Background(), tc.RuleType, tc.Opts...)
			if !errors.Is(err, tc.ExpectedError) {
				t.Fatalf("Expected error to be %v, got %v", tc.ExpectedError, err)
			}

			if len(rs) != tc.Matches {
				t.Fatalf("Expected %d matches, got %d", tc.Matches, len(rs))
			}

			for _, r := range rs {
				if r.Type != tc.RuleType {
					t.Fatalf("Expected all matched rules to have type %q, got %q", tc.RuleType, r.Type)
				}
			}
		})
	}
}
