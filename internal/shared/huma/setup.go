package huma

import (
	"net/http"

	"starter-boilerplate/internal/shared/config"

	gohuma "github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"go.uber.org/zap"
)

func Setup(mux *http.ServeMux, cfg config.AppConfig) gohuma.API {
	gohuma.NewErrorWithContext = func(_ gohuma.Context, status int, msg string, errs ...error) gohuma.StatusError {
		if status >= 500 {
			zap.L().Error("internal error", zap.Int("status", status), zap.Errors("errors", errs))
			return gohuma.NewError(status, "internal server error")
		}
		return gohuma.NewError(status, msg, errs...)
	}

	humaConfig := gohuma.DefaultConfig("Starter API", "1.0.0")

	if !cfg.SwaggerDocs {
		humaConfig.DocsPath = ""
		humaConfig.OpenAPIPath = ""
		humaConfig.SchemasPath = ""
	}

	api := humago.New(mux, humaConfig)

	if api.OpenAPI().Components.SecuritySchemes == nil {
		api.OpenAPI().Components.SecuritySchemes = make(map[string]*gohuma.SecurityScheme)
	}
	api.OpenAPI().Components.SecuritySchemes["bearerAuth"] = &gohuma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}

	return api
}
