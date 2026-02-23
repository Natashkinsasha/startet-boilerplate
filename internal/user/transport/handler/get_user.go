package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/transport/dto"

	"github.com/danielgtaylor/huma/v2"
)

type getUserInput struct {
	ID string `path:"id"`
}

type getUserOutput struct {
	Body struct {
		User dto.UserDTO `json:"user"`
	}
}

type GetUserHandler struct {
	uc *usecase.GetUserUseCase
}

func NewGetUserHandler(uc *usecase.GetUserUseCase) *GetUserHandler {
	return &GetUserHandler{uc: uc}
}

func (h *GetUserHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/v1/users/{id}",
		Summary:     "Get user by ID",
		Tags:        []string{"users"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, h.handle)
}

func (h *GetUserHandler) handle(ctx context.Context, input *getUserInput) (*getUserOutput, error) {
	u, err := h.uc.Execute(middleware.NewUserCtx(ctx), input.ID)
	if err != nil {
		return nil, err
	}

	out := &getUserOutput{}
	out.Body.User = dto.NewUserDTO(u)
	return out, nil
}
