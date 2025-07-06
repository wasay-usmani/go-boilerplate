package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const (
	AppName = "go-boilerplate"
	GitRepo = "go-boilerplate"
)

type Config struct {
	ListenHost    string `validate:"required"`
	ListenPort    string `validate:"required"`
	RPCListenPort string `validate:"required"`
	LogLevel      string `validate:"required"`
	Environment   string `validate:"required"`
	Debug         bool

	DBSchema             string `validate:"required"`
	SuperUserDatabaseURL string `validate:"required"`
	WriteDBURL           string `validate:"required"`
	ReadDBURL            string `validate:"required"`
}

// LoadConfig returns app configuration, recommend setting app build as
// build arg for proper/expected log output
// LoadConfig reads viper.Viper config using gitrepo & appName in the following paths
//
//	config
//	/etc/#{appName}/config
//	./internal/#{appName}/config
//	../config
//	../../config
//	../../../../#{gitRepo}/internal/#{appName}/config
//	../../../#{gitRepo}/internal/#{appName}/config
func LoadConfig(appBuild string) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath("config")
	v.AddConfigPath(fmt.Sprintf("/etc/%s/config", AppName))
	v.AddConfigPath(fmt.Sprintf("./internal/%s/config", AppName))
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath(fmt.Sprintf("../../../../%s/internal/%s/config", GitRepo, AppName))
	v.AddConfigPath(fmt.Sprintf("../../../%s/internal/%s/config", GitRepo, AppName))

	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using ENV vars. Note: This is expected behavior in k8s. Msg:	 %s\n", err)
	}

	// Determine Runtime Environment
	config := Config{
		Debug:                v.GetBool("DEBUG"),
		LogLevel:             v.GetString("LOG_LEVEL"),
		Environment:          v.GetString("ENVIRONMENT"),
		ListenHost:           v.GetString("HTTP_LISTEN_HOST"),
		ListenPort:           v.GetString("HTTP_LISTEN_PORT"),
		DBSchema:             v.GetString("DB_SCHEMA"),
		SuperUserDatabaseURL: v.GetString("SUPERUSER_DATABASE_URL"),
		WriteDBURL:           v.GetString("WRITE_DB_URL"),
		ReadDBURL:            v.GetString("READ_DB_URL"),
	}

	// Validate Config
	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
