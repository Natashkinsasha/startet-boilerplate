package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/transport/dto"

	"github.com/danielgtaylor/huma/v2"
)

type registerInput struct {
	Body struct {
		Email    string `json:"email" required:"true" format:"email"`
		Password string `json:"password" required:"true" minLength:"6"`
	}
}

type RegisterHandler struct {
	registerUC *usecase.RegisterUseCase
}

func NewRegisterHandler(registerUC *usecase.RegisterUseCase) *RegisterHandler {
	return &RegisterHandler{registerUC: registerUC}
}

func (h *RegisterHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-register",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/register",
		Summary:     "Register",
		Tags:        []string{"auth"},
	}, h.handle)
}

func (h *RegisterHandler) handle(ctx context.Context, input *registerInput) (*tokenOutput, error) {
	pair, err := h.registerUC.Execute(ctx, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, err
	}

	return &tokenOutput{Body: dto.NewTokenPairDTO(pair)}, nil
}
