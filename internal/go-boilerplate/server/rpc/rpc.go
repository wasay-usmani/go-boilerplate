package rpc

import (
	"fmt"
	"net"

	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/app"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/config"
	"google.golang.org/grpc"
)

type H struct {
	conf   *config.Config
	server *grpc.Server
	a      *app.Module
}

// Creates a new rpc handler
func NewHandlerBase(c *config.Config, application *app.Module) *H {
	handler := &H{a: application, conf: c}
	handler.server = grpc.NewServer()
	return handler
}

func (h *H) Run() error {
	listener, err := net.Listen("tcp", h.conf.RPCListenPort)
	if err != nil {
		return fmt.Errorf("failed to initialize rpc listener: %s", err.Error())
	}

	if err := h.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %s", err.Error())
	}

	return nil
}

func (h *H) Stop() {
	h.server.Stop()
}
