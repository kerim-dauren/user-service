package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	AppName  string    `env:"APP_NAME" env-default:"user-service"`
	AppEnv   string    `env:"APP_ENV" env-default:"dev"`
	HttpPort int       `env:"HTTP_PORT" env-default:"8080"`
	GRPCPort int       `env:"GRPC_PORT" env-default:"8081"`
	DbUrl    string    `env:"DB_URL"`
	Log      LogConfig `env-prefix:"LOG_" env-default:"info"`
}

type LogConfig struct {
	Level   string `env:"LEVEL" env-default:"info"`
	Handler string `env:"HANDLER" env-default:"text"`
	Writer  string `env:"WRITER" env-default:"stdout"`
}

// LoadConfig reads configuration from a .env file (if it exists) and environment variables.
func LoadConfig() (Config, error) {
	var cfg Config
	_, err := os.Stat(".env")
	if err == nil {
		err = cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			return cfg, err
		}
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
