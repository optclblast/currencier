// Package app configures and runs application.
package app

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/optclblast/currencier/internal/config"
	v1 "github.com/optclblast/currencier/internal/controller/http/v1"
	"github.com/optclblast/currencier/internal/interface/httpserver"
	"github.com/optclblast/currencier/internal/pkg/logger"
	"github.com/optclblast/currencier/internal/usecase/cache"
	"github.com/optclblast/currencier/internal/usecase/interactor"
	"github.com/redis/go-redis/v9"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.NewLogger(logger.WithLevel(logger.MapLevel(cfg.Common.Level)))

	l.Debug("config", slog.Any("struct", cfg))

	cache := cache.NewCache(redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.URL,
		Username: cfg.Cache.User,
		Password: cfg.Cache.Secret,
	}))

	apiKey := os.Getenv("FXRATEAPI_API_TOKEN")

	currencyInteractor := interactor.NewCurrencyInteractor(
		l.WithGroup("currency-interactor"),
		apiKey,
		cache,
		http.DefaultClient,
	)

	handler := v1.NewHandler(l, v1.NewCurrencyController(
		l.WithGroup("currency-controller"),
		currencyInteractor,
	))

	currencyInteractor.RunWorker()

	httpServer := httpserver.New(l, handler, httpserver.Port(cfg.Rest.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(
			"interrupt signal. ",
			slog.String("signal", s.String()),
		)
	case err := <-httpServer.Notify():
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
