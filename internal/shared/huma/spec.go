package huma

import (
	"encoding/json"
	"os"
	"path/filepath"

	gohuma "github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

func GenerateSpecFile(api gohuma.API) {
	const path = "docs/swagger.json"

	spec, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		zap.L().Error("failed to marshal openapi spec", zap.Error(err))
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		zap.L().Error("failed to create docs dir", zap.Error(err))
		return
	}

	if err := os.WriteFile(path, spec, 0o644); err != nil {
		zap.L().Error("failed to write swagger spec", zap.Error(err))
		return
	}

	zap.L().Info("swagger spec written", zap.String("path", path))
}
