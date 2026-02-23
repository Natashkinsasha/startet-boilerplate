package event

const PasswordChanged = "user.password_changed"

type PasswordChangedEvent struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

func (PasswordChangedEvent) EventName() string { return PasswordChanged }
func (PasswordChangedEvent) Tags() []string    { return []string{"profile"} }
