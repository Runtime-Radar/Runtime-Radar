package service

const (
	defaultPageSize = 10
	defaultOrder    = "created_at desc"
)

// gRPC errdetails.ErrorInfo.Reason codes used in service responses.
const (
	NameMustBeUnique = "NAME_MUST_BE_UNIQUE"
)
