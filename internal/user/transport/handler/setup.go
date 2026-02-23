package handler

import (
	"github.com/danielgtaylor/huma/v2"
)

type HandlersInit struct{}

func SetupHandlers(api huma.API, loginH *LoginHandler, refreshH *RefreshHandler, getUserH *GetUserHandler, registerH *RegisterHandler, changePasswordH *ChangePasswordHandler) HandlersInit {
	loginH.Register(api)
	refreshH.Register(api)
	getUserH.Register(api)
	registerH.Register(api)
	changePasswordH.Register(api)
	return HandlersInit{}
}
