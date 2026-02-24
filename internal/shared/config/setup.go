package config

import (
	"fmt"
	"os"
	"time"

	pkgamqp "starter-boilerplate/pkg/amqp"
	pkgdb "starter-boilerplate/pkg/db"
	pkggrpc "starter-boilerplate/pkg/grpc"
	"starter-boilerplate/pkg/outbox"
	pkgredis "starter-boilerplate/pkg/redis"

	sharedjwt "starter-boilerplate/internal/shared/jwt"
	sharedlogger "starter-boilerplate/internal/shared/logger"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port            int           `yaml:"port" validate:"required"`
	ReadTimeout     time.Duration `yaml:"read_timeout" validate:"required"`
	WriteTimeout    time.Duration `yaml:"write_timeout" validate:"required"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" validate:"required"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" validate:"required"`
	SwaggerDocs     bool          `yaml:"swagger_docs"`
	SwaggerFile     bool          `yaml:"swagger_file"`
}

type Config struct {
	App    AppConfig                 `yaml:"app"`
	Logger sharedlogger.LoggerConfig `yaml:"logger"`
	DB     pkgdb.DBConfig            `yaml:"db"`
	Redis  pkgredis.RedisConfig      `yaml:"redis"`
	JWT    sharedjwt.JWTConfig       `yaml:"jwt"`
	GRPC   pkggrpc.GRPCConfig        `yaml:"grpc"`
	AMQP   pkgamqp.AMQPConfig        `yaml:"amqp"`
	Outbox outbox.RelayConfig        `yaml:"outbox"`
}

func SetupConfig() *Config {
	cfg := &Config{}
	readYAML("env/.env.yaml", cfg)

	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		readYAML(fmt.Sprintf("env/.env.%s.yaml", appEnv), cfg)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic("invalid config: " + err.Error())
	}
	return cfg
}

func readYAML(path string, cfg *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read config file " + path + ": " + err.Error())
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		panic("failed to parse config file " + path + ": " + err.Error())
	}
}
