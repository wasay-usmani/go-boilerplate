package configkit

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type HTTP struct {
	ListenHost      string `validate:"required"`
	ListenPort      string `validate:"required"`
	Timeout         int
	EnableKeepAlive bool
}

func (c *C) LoadHTTPConfig() (*HTTP, error) {
	cfg := &HTTP{
		ListenHost:      c.Viper.GetString("http_listen_host"),
		ListenPort:      c.Viper.GetString("http_listen_port"),
		Timeout:         c.Viper.GetInt("http_timeout"),
		EnableKeepAlive: c.Viper.GetBool("http_enable_keep_alive"),
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("http config validation failed: %w", err)
	}

	return cfg, nil
}
