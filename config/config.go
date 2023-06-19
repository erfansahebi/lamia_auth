package config

import (
	"github.com/erfansahebi/lamia_shared/services"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"path"
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

	JWT struct {
		Secret   string `env:"JWT_SECRET"`
		Duration uint   `env:"JWT_EXPIRE_DURATION_MINUTE"`
	}
}

func LoadConfig() (*Config, error) {
	var configuration Config

	configPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err = godotenv.Load(path.Join(configPath, ".env")); err != nil {
		return nil, err
	}

	if err = cleanenv.ReadEnv(&configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}
