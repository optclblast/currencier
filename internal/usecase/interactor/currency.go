package interactor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/optclblast/currencier/internal/entity"
	"github.com/optclblast/currencier/internal/pkg/logger"
	"github.com/optclblast/currencier/internal/usecase/cache"
	"github.com/optclblast/currencier/internal/usecase/repo"
)

type CurrencyInteractor interface {
	GetQuotation(
		ctx context.Context,
		params GetQuotationParams,
	) (*entity.CurrencyQuotation, error)
}

type currencyInteractor struct {
	log    *slog.Logger
	repo   repo.QutesRepo
	cache  cache.Cache
	client *http.Client
	apiKey string
}

func NewCurrencyInteractor(
	log *slog.Logger,
	//repo repo.QutesRepo,
	cache cache.Cache,
	client *http.Client,
) CurrencyInteractor {
	return &currencyInteractor{
		log: log,
		//repo:   repo,
		cache:  cache,
		client: client,
	}
}

type GetQuotationParams struct {
	Date         time.Time
	CurrencyID   string
	CurrencyIDTo string
}

type currencyApiResponse struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Base    string             `json:"base"`
}

const baseCurrURL string = "https://api.fxratesapi.com/historical"

func (c *currencyInteractor) GetQuotation(
	ctx context.Context,
	params GetQuotationParams,
) (*entity.CurrencyQuotation, error) {
	if quotation, err := c.cache.Get(ctx, cache.GetParams{
		CurrencyID:   params.CurrencyID,
		CurrencyIDTo: params.CurrencyIDTo,
		DateKey:      buildDateKey(params.Date),
	}); err == nil {
		return quotation, nil
	}

	if params.CurrencyIDTo == "" {
		params.CurrencyIDTo = "RUB"
	}

	cdrUrl, err := c.requestCurrencyQuotationsURL(params.Date, params.CurrencyIDTo, params.CurrencyID)
	if err != nil {
		return nil, fmt.Errorf("error build request url. %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		cdrUrl,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("error prepare request. %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error do request to %s. %w", cdrUrl, err)
	}

	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read cdr response body. %w", err)
	}

	data := new(currencyApiResponse)

	if err := json.Unmarshal(rawBody, data); err != nil {
		return nil, fmt.Errorf("error unmarshal response body. %w", err)
	}

	quote := entity.CurrencyQuotation{
		Date:         params.Date,
		CurrencyID:   params.CurrencyID,
		CurrencyIDTo: params.CurrencyIDTo,
		Value:        data.Rates[params.CurrencyIDTo],
	}

	if err = c.cache.Set(ctx, quote, time.Hour*24); err != nil {
		c.log.Error(
			"error update cache",
			logger.Err(fmt.Errorf("error update cache. %w", err)),
		)
	}

	return &quote, nil
}

func (c *currencyInteractor) requestCurrencyQuotationsURL(
	t time.Time,
	currTo string,
	currFrom string,
) (string, error) {
	u, err := url.Parse(baseCurrURL)
	if err != nil {
		return "", fmt.Errorf("error parse cdr base url. %w", err)
	}

	if t.IsZero() {
		t = time.Now()
	}

	q := u.Query()
	q.Set("api_key", c.apiKey)
	q.Set("date", buildDateKey(t))
	q.Set("places", "6")
	q.Set("format", "json")
	q.Set("currencies", currTo)
	q.Set("base", currFrom)

	return u.String() + "?" + q.Encode(), nil
}

func buildDateKey(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}

	return fmt.Sprintf("&date=%v-%v-%v", t.Year(), int(t.Month()), t.Day())
}
