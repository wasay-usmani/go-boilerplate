package configkit

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// C holds the application configuration.
type C struct {
	Viper *viper.Viper

	AppName     string `validate:"required"`
	AppBuild    string `validate:"required"`
	LogLevel    string `validate:"required"`
	Environment string `validate:"required"`
	Debug       bool   `validate:"required"`
}

func LoadCacheConfig(appBuild string) (*C, error) {
	v := viper.New()
	v.SetConfigType("json")
	v.AutomaticEnv()

	cfg := &C{Viper: v, AppBuild: appBuild}

	// Populate the typed Config struct fields
	cfg.AppName = v.GetString("app_name")
	cfg.LogLevel = v.GetString("log_level")
	cfg.Environment = v.GetString("environment")
	cfg.Debug = v.GetBool("debug")

	// Validate required fields
	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
