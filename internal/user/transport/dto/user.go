package dto

import "starter-boilerplate/internal/user/domain/model"

type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TokenPairDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserDTO(u *model.User) UserDTO {
	return UserDTO{
		ID:    u.ID,
		Email: u.Email,
		Role:  string(u.Role),
	}
}

func NewTokenPairDTO(tp *model.TokenPair) TokenPairDTO {
	return TokenPairDTO{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
	}
}
