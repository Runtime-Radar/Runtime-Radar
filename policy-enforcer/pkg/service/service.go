package service

const (
	defaultPageSize = 10
	defaultOrder    = "created_at desc"
)

// gRPC errdetails.ErrorInfo.Reason codes used in service responses.
const (
	NameMustBeUnique         = "NAME_MUST_BE_UNIQUE"
	PredeclaredScopeDeletion = "PREDECLARED_SCOPE_DELETION"
	HasLinkedUsers           = "HAS_LINKED_USERS"
	HasLinkedRules           = "HAS_LINKED_RULES"
	// indicates that rules that were passed as enforced rules don't exist in database
	MissingEnforcedRules = "MISSING_ENFORCED_RULES"
)
