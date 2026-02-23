package domain

const EventUserCreated = "user.created"

type UserCreatedEvent struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Email  string `json:"email"   validate:"required,email"`
}

func (UserCreatedEvent) EventName() string { return EventUserCreated }
