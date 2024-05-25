// Package app configures and runs application.
package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/optclblast/currencier/internal/config"
	v1 "github.com/optclblast/currencier/internal/controller/http/v1"
	"github.com/optclblast/currencier/internal/interface/httpserver"
	"github.com/optclblast/currencier/internal/pkg/logger"
	"github.com/optclblast/currencier/internal/pkg/postgres"
	// "github.com/optclblast/currencier/internal/usecase/webapi"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.NewLogger(logger.MapLevel(cfg.Common.Level))

	l.Debug("config", slog.Any("struct", cfg))

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Error(
			"error initialize pg connection",
			logger.Err(err),
		)
	}
	defer pg.Close()

	// Use case
	// translationUseCase := usecase.New(
	// 	repo.New(pg),
	// 	webapi.New(),
	// )

	// HTTP Server
	handler := v1.NewHandler(l)
	httpServer := httpserver.New(l, handler, httpserver.Port(cfg.Rest.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(
			"interrupt signal. ",
			slog.String("signal", s.String()),
		)
	case err = <-httpServer.Notify():
		l.Error(
			"error",
			logger.Err(err),
		)

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			l.Error(
				"error shut down server",
				logger.Err(err),
			)
		}
	}
}
