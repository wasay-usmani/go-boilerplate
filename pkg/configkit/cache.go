package configkit

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/wasay-usmani/go-boilerplate/pkg/utils"
)

type Cache struct {
	Addresses  []string `validate:"required"`
	TLSEnabled *bool    `validate:"required"`
	KeySpace   string   `validate:"required"`
	Username   string
	Password   string
	ClientName string
}

func (c *C) LoadConfig(prefix string) (*Cache, error) {
	cfg := Cache{
		Addresses:  c.Viper.GetStringSlice("prefix_cache_addresses"),
		Username:   c.Viper.GetString("prefix_cache_username"),
		Password:   c.Viper.GetString("prefix_cache_password"),
		ClientName: c.Viper.GetString("prefix_cache_client_name"),
		TLSEnabled: utils.PtrOf(c.Viper.GetBool("prefix_cache_tls_enabled")),
		KeySpace:   c.Viper.GetString("prefix_cache_keyspace"),
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("cache config validation failed: %w", err)
	}

	return &cfg, nil
}
