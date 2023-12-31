// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"go-dk/messaging"
	"go-dk/storage"
)

type Server struct {
	address  string
	database *storage.Database
	log      *zap.Logger
	mux      chi.Router
	queue    *messaging.Queue
	server   *http.Server
}

type Options struct {
	Database *storage.Database
	Host     string
	Port     int
	Log      *zap.Logger
	Queue    *messaging.Queue
}

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()
	return &Server{
		address:  address,
		database: opts.Database,
		log:      opts.Log,
		mux:      mux,
		queue:    opts.Queue,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the server by setting routes and listening on the given address.
func (s *Server) Start() error {
	if err := s.database.Connect(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	s.setupRoutes()

	s.log.Info("Starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Stop the server gracefully within the timeout
func (s *Server) Stop() error {
	s.log.Info("Stopping")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping the server: %w", err)
	}

	return nil
}
