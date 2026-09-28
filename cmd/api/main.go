package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := loadConfig()
	app := &application{
		config:       cfg,
		logger:       slog.New(slog.NewTextHandler(os.Stdout, nil)),
		users:        newUserStore(),
		accounts:     newAccountStore(),
		categories:   newCategoryStore(),
		transactions: newTransactionStore(),
	}
	server := &http.Server{
		Addr:              app.config.httpAddr,
		Handler:           app.router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	app.logger.Info("starting server", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}
