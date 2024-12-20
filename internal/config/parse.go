package config

import (
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"
	"go.uber.org/zap"

	"github.com/Slava02/ChatSupport/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't open config file: %w", err)
	}
	defer file.Close()

	var config Config

	b, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't read config file: %w", err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't unmarshall config file: %w", err)
	}

	err = validator.Validator.Struct(config)
	if err != nil {
		return Config{}, fmt.Errorf("validation error: %w", err)
	}

	if config.Clients.Keycloak.DebugMode && config.Global.Env == "prod" {
		zap.L().Warn("Keycloak is set to debug mode while app is in production mode")
	}

	return config, nil
}
