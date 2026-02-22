package jwt

import (
	"time"

	pkgjwt "starter-boilerplate/pkg/jwt"
)

type JWTConfig struct {
	AccessSecret  string        `yaml:"access_secret" validate:"required"`
	RefreshSecret string        `yaml:"refresh_secret" validate:"required"`
	AccessTTL     time.Duration `yaml:"access_ttl" validate:"required"`
	RefreshTTL    time.Duration `yaml:"refresh_ttl" validate:"required"`
}

func NewJWTManager(cfg JWTConfig) *pkgjwt.Manager {
	return pkgjwt.NewManager(pkgjwt.Config{
		AccessSecret:  cfg.AccessSecret,
		RefreshSecret: cfg.RefreshSecret,
		AccessTTL:     cfg.AccessTTL,
		RefreshTTL:    cfg.RefreshTTL,
	})
}
