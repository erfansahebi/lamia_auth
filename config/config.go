package config

import (
	"github.com/erfansahebi/lamia_shared/services"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	services.Config

	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Username string `env:"REDIS_USERNAME"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB"`
	}

	AuthorizationToken struct {
		Duration uint `env:"AUTHORIZATION_TOKEN_EXPIRE_DURATION_MINUTE"`
	}
}

func LoadConfig() (*Config, error) {
	var configuration Config

	if err := cleanenv.ReadEnv(&configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}
