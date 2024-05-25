// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

// Http server
type Server struct {
	server          *http.Server
	log             *slog.Logger
	notify          chan error
	shutdownTimeout time.Duration
}

// New creates a new http server
func New(log *slog.Logger, handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		log:             log,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	s.log.Info(
		"starting http server",
		slog.String("host", s.server.Addr),
	)

	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify return error channel.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown begins a graceful shutdown.
func (s *Server) Shutdown() error {
	s.log.Info("shutting down http server. bye bye!")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
