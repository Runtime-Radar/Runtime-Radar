package service

import enf_api "github.com/runtime-radar/runtime-radar/policy-enforcer/api"

type Rule interface {
	enf_api.RuleControllerClient
}
