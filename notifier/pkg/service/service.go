package service

const (
	defaultOrder = "created_at desc"
)

// gRPC errdetails.ErrorInfo.Reason codes used in service responses.
const (
	NameMustBeUnique        = "NAME_MUST_BE_UNIQUE"
	IntegrationInaccessible = "INTEGRATION_INACCESSIBLE"
	NotificationInUse       = "NOTIFICATION_IN_USE"
	NotificationFailed      = "NOTIFICATION_FAILED"
)
