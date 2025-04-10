package goboilerplate

import "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/repository"

type QrysHandler interface {
}

type qrysHandler struct {
	r *repository.Module
}

func NewQrysHandler(r *repository.Module) QrysHandler {
	return &qrysHandler{
		r: r,
	}
}
