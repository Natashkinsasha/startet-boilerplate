package huma

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	gohuma "github.com/danielgtaylor/huma/v2"
)

func GenerateSpecFile(api gohuma.API) {
	const path = "docs/swagger.json"

	spec, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		slog.Error("failed to marshal openapi spec", slog.Any("error", err))
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		slog.Error("failed to create docs dir", slog.Any("error", err))
		return
	}

	if err := os.WriteFile(path, spec, 0o644); err != nil {
		slog.Error("failed to write swagger spec", slog.Any("error", err))
		return
	}

	slog.Info("swagger spec written", slog.String("path", path))
}
