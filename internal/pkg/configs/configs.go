package configs

import (
	"doctormakarhina/lumos/internal/pkg/errs"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const DotEnvEnabledEnvName = "DOTENV_ENABLED"
const DotEnvPathEnvName = "DOTENV_CONFIG_PATH"
const DefaultDotEnvPath = ".env"

type Validator interface {
	Validate() error
}

func LoadDotenvIfEnabled() {
	dotEnvEnabled := os.Getenv(DotEnvEnabledEnvName)
	if dotEnvEnabled == "" || dotEnvEnabled == "false" || dotEnvEnabled == "0" {
		return
	}
	configPath := os.Getenv(DotEnvPathEnvName)
	if configPath == "" {
		configPath = DefaultDotEnvPath
	}
	err := godotenv.Load(configPath)
	if err != nil {
		slog.Warn(
			fmt.Sprintf("Failed to load .env file %s", configPath),
			slog.String("err", err.Error()),
		)
	}
}

func Load(conf any) error {
	if err := env.Parse(conf); err != nil {
		return errs.WrapErrorf(err, errs.ErrCodeParsingFailed, "failed to parse env config")
	}

	if validator, ok := conf.(Validator); ok {
		err := validator.Validate()
		if err != nil {
			return errs.WrapErrorf(err, errs.ErrCodeInvalidArgument, "env config validation failed")
		}
	}

	return nil
}
