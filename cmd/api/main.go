package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	serverErrors := make(chan error, 1)

	go func() {
		app.logger.Info("starting server", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()
	shutdownSignal := make(chan os.Signal, 1)

	signal.Notify(shutdownSignal, os.Interrupt)
	defer signal.Stop(shutdownSignal)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			app.logger.Error("server stopped with error", "error", err)
			os.Exit(1)
		}
	case receivedSignal := <-shutdownSignal:
		app.logger.Info(
			"shutdown signal received",
			"signal",
			receivedSignal,
		)
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		app.logger.Error(
			"failed to shutdown server gracefully",
			"error",
			err,
		)
		os.Exit(1)
	}
	app.logger.Info("server stopped")
}
