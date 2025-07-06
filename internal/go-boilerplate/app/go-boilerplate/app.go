package goboilerplate

import "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/repository"

type App interface {
	CmdsHandler
	QrysHandler
}

type app struct {
	Cmds CmdsHandler
	Qrys QrysHandler
}

func New(writeDir, readDir *repository.Module) App {
	return &app{
		Cmds: NewCmdsHandler(writeDir),
		Qrys: NewQrysHandler(readDir),
	}
}
