package v1

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/optclblast/currencier/internal/usecase/interactor"
)

type CurrencyController interface {
	GetCurrencyQuotation(w http.ResponseWriter, req *http.Request) (any, error)
}

type currencyController struct {
	log        *slog.Logger
	interactor interactor.CurrencyInteractor
}

func NewCurrencyController(
	log *slog.Logger,
	interactor interactor.CurrencyInteractor,
) CurrencyController {
	return &currencyController{
		log:        log,
		interactor: interactor,
	}
}

func (c *currencyController) GetCurrencyQuotation(w http.ResponseWriter, req *http.Request) (any, error) {
	queryParams := req.URL.Query()

	if _, ok := queryParams["date"]; !ok {
		return nil, ErrorDateRequired
	}

	if _, ok := queryParams["val"]; !ok {
		return nil, ErrorValRequired
	}

	layout := "02.01.2006"

	date, err := time.Parse(layout, queryParams["date"][0])
	if err != nil {
		return nil, ErrorDateInvalid
	}

	var currTo string = "RUB"

	if valTo, ok := queryParams["val_to"]; ok {
		currTo = valTo[0]
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*5)
	defer cancel()

	q, err := c.interactor.GetQuotation(ctx, interactor.GetQuotationParams{
		Date:         date,
		CurrencyID:   queryParams["val"][0],
		CurrencyIDTo: currTo,
	})
	if err != nil {
		return nil, err
	}

	return q, nil
}
