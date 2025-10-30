package jwt

type Action string

const (
	ActionCreate  Action = "create"
	ActionRead    Action = "read"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionExecute Action = "execute"
)

type Actions []Action

func (as Actions) contains(v Action) bool {
	for _, a := range as {
		if a == v {
			return true
		}
	}
	return false
}

func (as Actions) containsAll(vs ...Action) bool {
	for _, a := range vs {
		if !as.contains(a) {
			return false
		}
	}
	return true
}
