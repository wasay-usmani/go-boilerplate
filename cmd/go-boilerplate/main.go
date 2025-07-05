package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/app"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/config"
	http_server "github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/server/http"
	"github.com/wasay-usmani/go-boilerplate/internal/go-boilerplate/server/rpc"
)

const _shutdownTimeout = 5

var build string // 0:8 GIT SHA injected at build time in Dockerfile

func main() {
	buildString := build
	if buildString == "" {
		buildString = "testing-unset"
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Load config
	cfg, err := config.LoadConfig(buildString)
	if err != nil {
		log.Fatalln("load config error", err)
	}

	// Initialize app module
	appModule, appCleanUp := app.NewModule(cfg)

	// Initialize requests handler
	hBase := http_server.NewHandlerBase(cfg, appModule)
	router := hBase.LoadRoutes()

	// Start API Server
	server := &http.Server{
		Addr:         cfg.ListenHost + ":" + cfg.ListenPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}

	// Start RPC Server
	rpcHandle := rpc.NewHandlerBase(cfg, appModule)

	go func() {
		// Start Serving Connections
		rpcErr := rpcHandle.Run()
		if rpcErr != nil {
			log.Fatal(
				"Server Error while trying to serve rpc",
				"rpc.port", cfg.RPCListenPort,
				"err", rpcErr,
			)
		}
	}()

	go func() {
		// Start Serving Connections
		httpErr := server.ListenAndServe()
		if httpErr != nil && !errors.Is(httpErr, http.ErrServerClosed) {
			log.Fatal(
				"Server Error while trying to serve http",
				"http.host", cfg.ListenHost,
				"http.port", cfg.ListenPort,
				"err", httpErr,
			)
		}
	}()

	// Listen for Shutdown Signal
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so no need to add it
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Blocking op; waiting for signal
	<-quit

	// Shutdown HTTP Server
	// Shutdown HTTP server
	subCtx, shutdownCancel := context.WithTimeout(ctx, _shutdownTimeout*time.Second)
	_ = server.Shutdown(subCtx)

	close(quit)
	appCleanUp()
	// Uncomment when using zap logger
	// _ = logger.Sync()
	cancel()
	shutdownCancel()
}
