package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type server struct {
	port          uint
	keyFile       *keyFile
	stateFilePath string
}

func (s *server) listenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleFetchState)
	mux.HandleFunc("POST /{$}", s.handleUpdateState)
	mux.HandleFunc("DELETE /{$}", s.handlePurgeState)

	srv := http.Server{
		Addr:    net.JoinHostPort("localhost", fmt.Sprint(s.port)),
		Handler: mux,
	}

	// gracefully shutdown on SIGINT
	shutdown := make(chan os.Signal, 1)
	shutdownComplete := make(chan error, 1)
	go func() {
		<-shutdown
		slog.Info("shutting down")
		err := srv.Shutdown(context.Background())
		shutdownComplete <- err
	}()
	signal.Notify(shutdown, os.Interrupt)

	// listen and serve
	slog.Info("listening", "addr", srv.Addr)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// wait for shutdown to complete
	err = <-shutdownComplete
	return err
}

func (s *server) handleFetchState(w http.ResponseWriter, r *http.Request) {
	slog.Info("fetching state", "headers", r.Header)

	err := fetchState(w, s.keyFile, s.stateFilePath)
	if err != nil {
		slog.Error("failed to fetch state", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *server) handleUpdateState(w http.ResponseWriter, r *http.Request) {
	slog.Info("updating state", "headers", r.Header)

	err := updateState(r.Body, s.keyFile, s.stateFilePath)
	if err != nil {
		slog.Error("failed to update state", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *server) handlePurgeState(w http.ResponseWriter, r *http.Request) {
	slog.Info("purging state", "headers", r.Header)

	err := os.Remove(s.stateFilePath)
	if err != nil {
		slog.Error("failed to purge state", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
