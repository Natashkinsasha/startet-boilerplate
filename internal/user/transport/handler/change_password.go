package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type changePasswordInput struct {
	Body struct {
		OldPassword string `json:"old_password" required:"true" minLength:"6"`
		NewPassword string `json:"new_password" required:"true" minLength:"6"`
	}
}

type ChangePasswordHandler struct {
	uc *usecase.ChangePasswordUseCase
}

func NewChangePasswordHandler(uc *usecase.ChangePasswordUseCase) *ChangePasswordHandler {
	return &ChangePasswordHandler{uc: uc}
}

func (h *ChangePasswordHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "auth-change-password",
		Method:        http.MethodPut,
		Path:          "/api/v1/auth/password",
		Summary:       "Change password",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusNoContent,
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.handle)
}

func (h *ChangePasswordHandler) handle(ctx context.Context, input *changePasswordInput) (*struct{}, error) {
	err := h.uc.Execute(middleware.NewAuthCtx(ctx), input.Body.OldPassword, input.Body.NewPassword)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
