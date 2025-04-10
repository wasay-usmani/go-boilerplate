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
	Viper *viper.Viper

	ListenHost    string `validate:"required"`
	ListenPort    string `validate:"required"`
	RpcListenHost string `validate:"required"`
	RpcListenPort string `validate:"required"`
	AppName       string `validate:"required"`
	AppBuild      string `validate:"required"`
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
func LoadConfig(appBuild string) (*Config, error) {
	// Load common configs
	cfg, err := loadConfig(GitRepo, AppName, appBuild)
	if err != nil {
		return nil, err
	}

	config := Config{
		ListenHost:           cfg.Viper.GetString("HTTP_LISTEN_HOST"),
		ListenPort:           cfg.Viper.GetString("HTTP_LISTEN_PORT"),
		DBSchema:             cfg.Viper.GetString("DB_SCHEMA"),
		SuperUserDatabaseURL: cfg.Viper.GetString("SUPERUSER_DATABASE_URL"),
		WriteDBURL:           cfg.Viper.GetString("WRITE_DB_URL"),
		ReadDBURL:            cfg.Viper.GetString("READ_DB_URL"),
	}

	// Validate Config
	if err := validator.New().Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// LoadConfig reads viper.Viper config using gitrepo & appName in the following paths
//
//	config
//	/etc/#{appName}/config
//	./internal/#{appName}/config
//	../config
//	../../config
//	../../../../#{gitRepo}/internal/#{appName}/config
//	../../../#{gitRepo}/internal/#{appName}/config
func loadConfig(gitRepo, appName, appBuild string) (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath("config")
	v.AddConfigPath(fmt.Sprintf("/etc/%s/config", appName))
	v.AddConfigPath(fmt.Sprintf("./internal/%s/config", appName))
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath(fmt.Sprintf("../../../../%s/internal/%s/config", gitRepo, appName))
	v.AddConfigPath(fmt.Sprintf("../../../%s/internal/%s/config", gitRepo, appName))

	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using ENV vars. Note: This is expected behavior in k8s. Msg:	 %s\n", err)
	}

	// Determine Runtime Environment

	c := Config{
		Viper:       v,
		AppName:     appName,
		AppBuild:    appBuild,
		LogLevel:    v.GetString("LOG_LEVEL"),
		Environment: v.GetString("ENVIRONMENT"),
		Debug:       v.GetBool("DEBUG"),
	}

	// Validate Config
	if err := validator.New().Struct(c); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &c, nil
}
