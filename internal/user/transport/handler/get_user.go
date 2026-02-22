package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/shared/middleware"
	"starter-boilerplate/internal/user/app/service"

	"github.com/danielgtaylor/huma/v2"
)

type getUserInput struct {
	ID string `path:"id"`
}

type GetUserBody struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type getUserOutput struct {
	Body GetUserBody
}

type GetUserHandler struct {
	userService service.UserService
}

func NewGetUserHandler(userService service.UserService) *GetUserHandler {
	return &GetUserHandler{userService: userService}
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
	claims := middleware.ClaimsFromContext(ctx)

	if claims.Role != "admin" && claims.UserID != input.ID {
		return nil, huma.Error403Forbidden("access denied")
	}

	u, err := h.userService.FindByID(ctx, input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to fetch user")
	}
	if u == nil {
		return nil, huma.Error404NotFound("user not found")
	}

	return &getUserOutput{Body: GetUserBody{
		ID:    u.ID,
		Email: u.Email,
		Role:  string(u.Role),
	}}, nil
}
