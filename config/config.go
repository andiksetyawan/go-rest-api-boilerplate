package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type configs struct {
	ServiceName        string `env:"SERVICE_NAME" yaml:"service_name" env-default:"svc-go-rest-api-boilerplate"`
	ServiceAddress     string `env:"SERVICE_ADDRESS" yaml:"service_address" env-default:":8080"`
	ServiceEnvironment string `env:"SERVICE_ENVIRONMENT" yaml:"service_environment" env-default:"production"`

	DbHost string `env:"DB_HOST" yaml:"db_host" env-default:"localhost"`
	DbPort string `env:"DB_PORT" yaml:"db_port" env-default:"5436"`
	DbUser string `env:"DB_USER" yaml:"db_user" env-default:"postgres"`
	DbPass string `env:"DB_PASSWORD" yaml:"db_password" env-default:"postgres"`
	DbName string `env:"DB_NAME" yaml:"db_name" env-default:"svc-go-rest-api-boilerplate"`

	OtelUptraceDsn string `env:"OTEL_UPTRACE_DSN" yaml:"otel_uptrace_dsn" env-default:"https://ojnMDvABsRBuUbQntWnbnQ@uptrace.dev/860"`
	//OtelOtlpCollectorUrl string `env:"OTEL_OTLP_COLLECTOR_URL" yaml:"otel_otlp_collector_url" env-default:"localhost:4317"`
	//OtelInsecOtlpColUrl  bool   `env:"OTEL_INSECURE_OTLP_COLLECTOR" yaml:"otel_insecure_otlp_collector" env-default:"true"`
}

var App configs

func Init() {
	err := cleanenv.ReadConfig(".env", &App)
	if err != nil {
		log.Info("failed to read .env file, setting config from default environment variables.")
		cleanenv.ReadEnv(&App)
	}
}
