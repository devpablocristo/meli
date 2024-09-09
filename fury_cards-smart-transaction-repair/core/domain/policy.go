package domain

// UserType validation by  Auth-Policy Service.
const (
	UserType = "user"
)

type PolicyOutput struct {
	IsAuthorized    bool
	RestrictsFailed []string
}
