package http

import (
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/app"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/config"
)

type H struct {
	cfg *config.Config
	app *app.Module
}

func NewHandlerBase(cfg *config.Config, appModule *app.Module) *H {
	return &H{
		cfg: cfg,
		app: appModule,
	}
}
