package config

import (
	"os"
	"time"

	errutils "github.com/Egor213/PVZ/pkg/errors"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PG   `yaml:"postgres"`
		JWT  `yaml:"jwt"`
	}

	App struct {
		Name    string `yaml:"name" env-required:"true"`
		Version string `yaml:"version" env-required:"true"`
	}

	HTTP struct {
		Address string `env-required:"true" env:"SERVER_ADDRESS"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	}

	PG struct {
		URL         string `env-required:"true" env:"POSTGRES_CONN"`
		MaxPoolSize int    `env-required:"true" env:"MAX_POOL_SIZE" yaml:"max_pool_size"`
	}

	JWT struct {
		SignKey  string        `env-required:"true" env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-required:"true"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load("infra/.env"); err != nil {
		log.WithError(err).Info(".env file not found, using system environment")
	}

	pathToConfig, ok := os.LookupEnv("APP_CONFIG_PATH")
	if !ok || pathToConfig == "" {
		log.WithField("env_var", "APP_CONFIG_PATH").
			Info("Config path is not set, using default")
		pathToConfig = "config/config.yaml"
	}

	if err := cleanenv.ReadConfig(pathToConfig, cfg); err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, errutils.WrapPathErr(err)
	}

	return cfg, nil
}
