package configs

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Environment         string `envconfig:"ENVIRONMENT" required:"true"`
	LogErrorRequestDump bool   `envconfig:"LOG_ERROR_REQUEST_DUMP" required:"true"`
	ServerBindAddress   string `envconfig:"SERVER_BIND_ADDRESS" required:"true"`
	DbAddress           string `envconfig:"DB_ADDRESS" required:"true"`
	DbUser              string `envconfig:"DB_USER" required:"true"`
	DbPassword          string `envconfig:"DB_PASSWORD" required:"true"`
	LogLevel            string `envconfig:"LOG_LEVEL" required:"true"`
	DbName              string `envconfig:"DB_NAME" required:"true"`
	SessionName         string `envconfig:"SESSION_NAME" required:"true"`
	SessionSecretKey    string `envconfig:"SESSION_SECRET_KEY" required:"true"`
}

func LoadConfig() EnvConfig {
	var config EnvConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalln(err)
	}
	return config
}
