package app

import (
	goboilerplate "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/app/go-boilerplate"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/config"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/repository"
)

type Module struct {
	Boilerplate goboilerplate.App
}

func NewModule(cfg *config.Config) (module *Module, cleanup func()) {
	// Initialize read & write database connections
	// writeDB, readDB := initDBConns(cfg)
	writeR := repository.NewModule()
	readR := repository.NewModule()
	return &Module{
			Boilerplate: goboilerplate.New(writeR, readR),
		}, func() {
			// Database cleanup will be implemented when database connections are added
		}
}

func NewMockModule() *Module {
	return &Module{}
}
