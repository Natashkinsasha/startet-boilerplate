package handler

import (
	"context"
	"net"
	"net/http"

	"starter-boilerplate/internal/user/app/usecase"
	"starter-boilerplate/internal/user/transport/dto"

	"github.com/danielgtaylor/huma/v2"
)

type loginInput struct {
	IP        string `json:"-"`
	UserAgent string `json:"-"`
	Body      struct {
		Email    string `json:"email" required:"true" format:"email"`
		Password string `json:"password" required:"true" minLength:"6"`
	}
}

func (i *loginInput) Resolve(ctx huma.Context) []error {
	if ip := ctx.Header("X-Forwarded-For"); ip != "" {
		i.IP = ip
	} else if ip := ctx.Header("X-Real-IP"); ip != "" {
		i.IP = ip
	} else {
		host, _, err := net.SplitHostPort(ctx.RemoteAddr())
		if err != nil {
			host = ctx.RemoteAddr()
		}
		i.IP = host
	}

	i.UserAgent = ctx.Header("User-Agent")
	return nil
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
	pair, err := h.loginUC.Execute(ctx, input.Body.Email, input.Body.Password, input.IP, input.UserAgent)
	if err != nil {
		return nil, err
	}

	return &tokenOutput{Body: dto.NewTokenPairDTO(pair)}, nil
}
