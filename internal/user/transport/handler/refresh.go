package handler

import (
	"context"
	"net/http"

	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/transport/dto"

	"github.com/danielgtaylor/huma/v2"
)

type refreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" required:"true"`
	}
}

type RefreshHandler struct {
	refreshUC *usecase.RefreshUseCase
}

func NewRefreshHandler(refreshUC *usecase.RefreshUseCase) *RefreshHandler {
	return &RefreshHandler{refreshUC: refreshUC}
}

func (h *RefreshHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-refresh",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/refresh",
		Summary:     "Refresh tokens",
		Tags:        []string{"auth"},
	}, h.handle)
}

func (h *RefreshHandler) handle(ctx context.Context, input *refreshInput) (*tokenOutput, error) {
	pair, err := h.refreshUC.Execute(ctx, input.Body.RefreshToken)
	if err != nil {
		return nil, huma.Error401Unauthorized(err.Error())
	}

	return &tokenOutput{Body: dto.NewTokenPairDTO(pair)}, nil
}
