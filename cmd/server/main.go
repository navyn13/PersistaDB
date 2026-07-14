package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/navyn13/PersistaDB/internal/config"
	"github.com/navyn13/PersistaDB/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.NewServer(cfg)
	srv.Start()

	//server shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server gracefully...")
	srv.Shutdown()
	slog.Info("Server stopped.")
}
