package event

const UserLoggedIn = "user.logged_in"

type UserLoggedInEvent struct {
	UserID    string `json:"user_id"    validate:"required,uuid"`
	IP        string `json:"ip"         validate:"required"`
	UserAgent string `json:"user_agent" validate:"required"`
}

func (UserLoggedInEvent) EventName() string { return UserLoggedIn }
