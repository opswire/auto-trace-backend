package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type (
	// Config -.
	Config struct {
		App     App  `yaml:"ads"`
		Http    HTTP `yaml:"http"`
		Log     Log  `yaml:"logger"`
		Pg      PG   `yaml:"postgres"`
		Yokassa Yokassa
		Kafka   Kafka
		Smtp    Smtp
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

	Yokassa struct {
		Url                       string `env:"YOKASSA_URL"`
		RedirectUrl               string `env:"YOKASSA_REDIRECT_URL"`
		Username                  string `env:"YOKASSA_USERNAME"`
		Password                  string `env:"YOKASSA_PASSWORD"`
		WebhookAllowedIpAddresses string `env:"YOKASSA_WEBHOOK_WHITELIST"`
	}

	Kafka struct {
		Dsn           string `env:"KAFKA_DSN"`
		PaymentsTopic string `env:"KAFKA_PAYMENTS_TOPIC"`
	}

	Smtp struct {
		Host string `env:"SMTP_HOST"`
		Port string `env:"SMTP_PORT"`
		From string `env:"SMTP_FROM_EMAIL"`
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
