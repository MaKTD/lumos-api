package envconf

import (
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	DefaultDotEnvPath    = ".env"
	DotEnvPathEnvName    = "DOTENV_CONFIG_PATH"
	DotEnvEnabledEnvName = "DOTENV_ENABLED"
)

type Validator interface {
	Validate() error
}

func LoadDotenvIfEnabled() error {
	dotEnvEnabled := os.Getenv(DotEnvEnabledEnvName)
	dotEnvEnabled = strings.TrimSpace(dotEnvEnabled)
	dotEnvEnabled = strings.ToLower(dotEnvEnabled)

	if dotEnvEnabled != "1" && dotEnvEnabled != "true" {
		return nil
	}

	configPathsStr := os.Getenv(DotEnvPathEnvName)
	configPathsStr = strings.TrimSpace(configPathsStr)

	var configPaths []string
	if configPathsStr == "" {
		configPaths = []string{DefaultDotEnvPath}
	} else {
		configPaths = strings.Split(configPathsStr, ",")
	}

	return godotenv.Load(configPaths...)
}

func Load(conf any) error {
	if err := env.Parse(conf); err != nil {
		return err
	}

	if validator, ok := conf.(Validator); ok {
		err := validator.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
