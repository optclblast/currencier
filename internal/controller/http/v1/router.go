package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/optclblast/currencier/internal/pkg/logger"
)

type handler struct {
	*chi.Mux
	log        *slog.Logger
	controller CurrencyController
}

func NewHandler(
	log *slog.Logger,
	controller CurrencyController,
) *handler {
	h := &handler{
		Mux:        chi.NewRouter(),
		log:        log,
		controller: controller,
	}

	h.buildRouter()

	return h
}

func (s *handler) buildRouter() {
	r := chi.NewRouter()

	r.Get("/currency", s.handle(s.controller.GetCurrencyQuotation, "get_currency"))
}

func (s *handler) handle(
	h func(w http.ResponseWriter, req *http.Request) (any, error),
	methodName string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := h(w, r)
		if err != nil {
			s.responseError(w, err, methodName)

			return
		}

		w.WriteHeader(http.StatusOK)

		out, err := json.Marshal(resp)
		if err != nil {
			s.responseError(w, err, methodName)

			return
		}

		if _, err = w.Write(out); err != nil {
			s.log.Error(
				"error write http response",
				slog.String("method_name", methodName),
				logger.Err(err),
			)
		}
	}
}

func (s *handler) responseError(w http.ResponseWriter, e error, methodName string) {
	s.log.Error("error handle request", logger.Err(e), slog.String("method_name", methodName))

	apiErr := mapError(e)

	out, err := json.Marshal(apiErr)
	if err != nil {
		s.log.Error("error marshal api error", logger.Err(err))

		return
	}

	w.WriteHeader(apiErr.Code)
	w.Write(out)
}

func (s *handler) handleMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
