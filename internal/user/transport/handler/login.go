package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/user/app/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type loginInput struct {
	Body struct {
		Email    string `json:"email" required:"true" format:"email"`
		Password string `json:"password" required:"true" minLength:"6"`
	}
}

type LoginHandler struct {
	loginUC *usecase.LoginUseCase
}

func NewLoginHandler(loginUC *usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{loginUC: loginUC}
}

func (h *LoginHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-login",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/login",
		Summary:     "Login",
		Tags:        []string{"auth"},
	}, h.handle)
}

func (h *LoginHandler) handle(ctx context.Context, input *loginInput) (*tokenOutput, error) {
	pair, err := h.loginUC.Execute(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, huma.Error401Unauthorized(err.Error())
	}

	return &tokenOutput{Body: TokenBody{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}}, nil
}
