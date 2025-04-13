package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type (
	// Config -.
	Config struct {
		App  App  `yaml:"ads"`
		Http HTTP `yaml:"http"`
		Log  Log  `yaml:"logger"`
		Pg   PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax   int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL       string `env-required:"true"                 env:"PG_URL"`
		ExposeURL string `env-required:"true"                 env:"PG_EXPOSE_URL"`
	}
)

// NewConfig returns ads config.
func NewConfig() *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	return cfg
}
