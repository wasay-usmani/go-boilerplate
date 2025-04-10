package app

import (
	goboilerplate "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/app/go-boilerplate"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/config"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/repository"
)

type Module struct {
	Boilerplate goboilerplate.App
}

func NewModule(cfg *config.Config) (*Module, func()) {
	// Initialize read & write database connections
	// writeDB, readDB := initDBConns(cfg)

	writeR := repository.NewModule()
	readR := repository.NewModule()

	return &Module{
			Boilerplate: goboilerplate.New(writeR, readR),
		}, func() {
			//	_ = writeDB.Close()
			//	_ = readDB.Close()
		}
}

func NewMockModule() *Module {
	return &Module{}
}
