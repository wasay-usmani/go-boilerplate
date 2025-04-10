package goboilerplate

import "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/repository"

type CmdsHandler interface {
}

type cmdsHandler struct {
	r *repository.Module
}

func NewCmdsHandler(r *repository.Module) CmdsHandler {
	return &cmdsHandler{
		r: r,
	}
}
