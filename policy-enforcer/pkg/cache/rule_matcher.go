package cache

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/database"
	"github.com/runtime-radar/runtime-radar/policy-enforcer/pkg/model"
	"gorm.io/gorm"
)

const (
	ruleCacheKeyPrefix  = "rule_cache"
	ruleCacheExpiration = 24 * time.Hour
)

var (
	ErrNotScopeable = errors.New("rule type is not scopeable")
	ErrNoOptions    = errors.New("no options given")
)

// MatchValue represents matchable value which should be explicitly set as used
// to be taken into account during matching.
type MatchValue struct {
	Value string
	Used  bool
}

// MatchArgs represents the group of arguments that are used to match rules.
type MatchArgs struct {
	ImageName MatchValue
	Registry  MatchValue
	Namespace MatchValue
	Cluster   MatchValue
	Pod       MatchValue
	Container MatchValue
	Node      MatchValue
}

type MatchOption func(*MatchArgs)

func WithImageName(n string) MatchOption {
	return func(ma *MatchArgs) {
		ma.ImageName = MatchValue{n, true}
	}
}

func WithRegistry(r string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Registry = MatchValue{r, true}
	}
}

func WithNamespace(ns string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Namespace = MatchValue{ns, true}
	}
}

func WithCluster(c string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Cluster = MatchValue{c, true}
	}
}

func WithPod(p string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Pod = MatchValue{p, true}
	}
}

func WithContainer(c string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Container = MatchValue{c, true}
	}
}

func WithNode(n string) MatchOption {
	return func(ma *MatchArgs) {
		ma.Node = MatchValue{n, true}
	}
}

func NewMatchArgs(opts ...MatchOption) (MatchArgs, error) {
	if len(opts) == 0 {
		return MatchArgs{}, ErrNoOptions
	}

	args := MatchArgs{}

	for _, o := range opts {
		o(&args)
	}

	return args, nil
}

type RuleMatcher interface {
	// MatchRules return rules whose scope and type match given ones.
	// Scope is matched based only on given MatchOptions. Other fields are ignored.
	// If no MatchOptions are passed, ErrNoOptions is returned.
	// It's guaranteed that returned slice doesn't contain duplicates.
	MatchRules(context.Context, model.RuleType, ...MatchOption) ([]*model.Rule, error)
	Invalidate(context.Context) error
}

// ScopePatterns is model.Scope with all glob patterns compiled and ready to be used for matching.
type ScopePatterns struct {
	ImageNames []glob.Glob
	Namespaces []glob.Glob
	Clusters   []glob.Glob
	Registries []glob.Glob
	Pods       []glob.Glob
	Containers []glob.Glob
	Nodes      []glob.Glob
}

// RuleMatcherData contain rules as well as their scopes with globes compiled.
// Note that RuleMatcherData only stores rules whose scope is defined, because others cannot be used for matching.
type RuleMatcherData struct {
	RuleData map[uuid.UUID]*model.Rule
	// MatchData stores rule's id as a key and its scope's compiled patterns as a value.
	MatchData map[uuid.UUID]*ScopePatterns
}

func NewRuleMatcherData() *RuleMatcherData {
	return &RuleMatcherData{
		RuleData:  map[uuid.UUID]*model.Rule{},
		MatchData: map[uuid.UUID]*ScopePatterns{},
	}
}

func NewScopePatterns(s *model.Scope) (*ScopePatterns, error) {
	sp := &ScopePatterns{}

	for _, img := range s.ImageNames {
		g, err := glob.Compile(img)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", img, err)
		}

		sp.ImageNames = append(sp.ImageNames, g)
	}

	for _, reg := range s.Registries {
		g, err := glob.Compile(reg)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", reg, err)
		}

		sp.Registries = append(sp.Registries, g)
	}

	for _, nms := range s.Namespaces {
		g, err := glob.Compile(nms)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", nms, err)
		}

		sp.Namespaces = append(sp.Namespaces, g)
	}

	for _, cls := range s.Clusters {
		g, err := glob.Compile(cls)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", cls, err)
		}

		sp.Clusters = append(sp.Clusters, g)
	}

	for _, p := range s.Pods {
		g, err := glob.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", p, err)
		}

		sp.Pods = append(sp.Pods, g)
	}

	for _, c := range s.Containers {
		g, err := glob.Compile(c)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", c, err)
		}

		sp.Containers = append(sp.Containers, g)
	}

	for _, n := range s.Nodes {
		g, err := glob.Compile(n)
		if err != nil {
			return nil, fmt.Errorf("can't compile '%s' pattern: %w", n, err)
		}

		sp.Nodes = append(sp.Nodes, g)
	}

	return sp, nil
}

type RuleCache struct {
	Cache          Cache
	RuleRepository database.RuleRepository
}

func (rc *RuleCache) Invalidate(ctx context.Context) error {
	for _, rt := range model.ScopeableRuleTypes {
		if err := rc.Cache.Del(ctx, RuleCacheKey(rt)); err != nil {
			return fmt.Errorf("can't invalidate rule cache for type %s: %w", rt, err)
		}
	}

	log.Debug().Msgf("Rule cache invalidated")
	return nil
}

func (rc *RuleCache) MatchRules(ctx context.Context, ruleType model.RuleType, opts ...MatchOption) ([]*model.Rule, error) {
	args, err := NewMatchArgs(opts...)
	if err != nil {
		return nil, fmt.Errorf("can't create match args: %w", err)
	}

	if !slices.Contains(model.ScopeableRuleTypes, ruleType) {
		return nil, fmt.Errorf("%w: '%s'", ErrNotScopeable, ruleType)
	}

	rmd, err := rc.getRuleMatcherData(ctx, ruleType)
	if err != nil {
		return nil, fmt.Errorf("can't get rule matcher data: %w", err)
	}

	rs := []*model.Rule{}

	for id, m := range rmd.MatchData {
		matchImg, err := rc.matchValue(ctx, m.ImageNames, args.ImageName)
		if err != nil {
			return nil, fmt.Errorf("can't match image name: %w", err)
		}

		matchReg, err := rc.matchValue(ctx, m.Registries, args.Registry)
		if err != nil {
			return nil, fmt.Errorf("can't match registry: %w", err)
		}

		matchNms, err := rc.matchValue(ctx, m.Namespaces, args.Namespace)
		if err != nil {
			return nil, fmt.Errorf("can't match namespace: %w", err)
		}

		matchCls, err := rc.matchValue(ctx, m.Clusters, args.Cluster)
		if err != nil {
			return nil, fmt.Errorf("can't match cluster: %w", err)
		}

		matchPods, err := rc.matchValue(ctx, m.Pods, args.Pod)
		if err != nil {
			return nil, fmt.Errorf("can't match pod: %w", err)
		}

		matchConts, err := rc.matchValue(ctx, m.Containers, args.Container)
		if err != nil {
			return nil, fmt.Errorf("can't match container: %w", err)
		}

		matchNodes, err := rc.matchValue(ctx, m.Nodes, args.Node)
		if err != nil {
			return nil, fmt.Errorf("can't match node: %w", err)
		}

		r := rmd.RuleData[id]

		if ruleType == r.Type &&
			matchImg &&
			matchReg &&
			matchNms &&
			matchCls &&
			matchPods &&
			matchConts &&
			matchNodes {
			rs = append(rs, r)
		}
	}

	return rs, nil
}

// matchValue returns true if given mv matches at least one of gs.
// If mv is not set as used explicitly, it's also considered as matching in order to simplify logics of MatchRules.
func (rc *RuleCache) matchValue(ctx context.Context, gs []glob.Glob, mv MatchValue) (bool, error) {
	if !mv.Used {
		return true, nil
	}

	for _, g := range gs {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			if g.Match(mv.Value) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (rc *RuleCache) getRuleMatcherData(ctx context.Context, rt model.RuleType) (*RuleMatcherData, error) {
	var rmd *RuleMatcherData

	key := RuleCacheKey(rt)

	ok, err := rc.Cache.Get(ctx, key, &rmd)
	if err != nil {
		return nil, fmt.Errorf("can't get rule matcher data from cache: %w", err)
	}
	if ok {
		return rmd, nil
	}

	filter := gorm.Expr("type = ? and scope is not null", rt)
	rs, err := rc.RuleRepository.GetAll(ctx, filter, nil, true) // preload is on
	if err != nil {
		return nil, fmt.Errorf("can't get all rules from repository: %w", err)
	}

	rmd, err = rc.populate(ctx, key, rs)
	if err != nil {
		return nil, fmt.Errorf("can't populate rule cache: %w", err)
	}

	return rmd, nil
}

func (rc *RuleCache) populate(ctx context.Context, key string, rs []*model.Rule) (*RuleMatcherData, error) {
	rmd := NewRuleMatcherData()

	for _, r := range rs {
		rmd.RuleData[r.ID] = r

		sm, err := NewScopePatterns(r.Scope)
		if err != nil {
			return nil, fmt.Errorf("can't build new scope matcher: %w", err)
		}

		rmd.MatchData[r.ID] = sm
	}

	if err := rc.Cache.Set(ctx, key, rmd, ruleCacheExpiration); err != nil {
		return nil, fmt.Errorf("can't save rule matcher data to cache: %w", err)
	}
	log.Debug().Msgf("Populated cache with %d RuleData and %d MatchData entries", len(rmd.RuleData), len(rmd.MatchData))

	return rmd, nil
}

func RuleCacheKey(rt model.RuleType) string {
	return ruleCacheKeyPrefix + "_" + rt.String()
}
