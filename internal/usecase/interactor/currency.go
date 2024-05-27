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
)

type CurrencyInteractor interface {
	GetQuotation(
		ctx context.Context,
		params GetQuotationParams,
	) (*entity.CurrencyQuotation, error)
	RunWorker()
}

type currencyInteractor struct {
	log    *slog.Logger
	cache  cache.Cache
	client *http.Client
	apiKey string
}

func NewCurrencyInteractor(
	log *slog.Logger,
	apiKey string,
	cache cache.Cache,
	client *http.Client,
) CurrencyInteractor {
	return &currencyInteractor{
		log:    log,
		apiKey: apiKey,
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
	Success     bool               `json:"success"`
	Rates       map[string]float64 `json:"rates,omitempty"`
	Base        string             `json:"base,omitempty"`
	Error       string             `json:"error,omitempty"`
	Description string             `json:"description,omitempty"`
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
		c.log.Debug(
			"cache hit",
			slog.String("CurrencyID", params.CurrencyID),
			slog.String("CurrencyIDTo", params.CurrencyIDTo),
			slog.String("DateKey", buildDateKey(params.Date)),
		)
		return quotation, nil
	} else {
		c.log.Debug(
			"cache miss",
			slog.String("CurrencyID", params.CurrencyID),
			slog.String("CurrencyIDTo", params.CurrencyIDTo),
			slog.String("DateKey", buildDateKey(params.Date)),
			logger.Err(err),
		)
	}

	if params.CurrencyIDTo == "" {
		params.CurrencyIDTo = "RUB"
	}

	cdrUrl, err := c.requestCurrencyQuotationsURL(params.Date, params.CurrencyIDTo, params.CurrencyID)
	if err != nil {
		return nil, fmt.Errorf("error build request url. %w", err)
	}

	quote, err := c.fetchQuote(ctx, cdrUrl, params)
	if err != nil {
		return nil, fmt.Errorf("error fetch quote. %w", err)
	}

	if err = c.cache.Set(
		ctx,
		params.CurrencyID+params.CurrencyIDTo+buildDateKey(params.Date),
		*quote,
		time.Hour*24,
	); err != nil {
		c.log.Error(
			"error update cache",
			logger.Err(fmt.Errorf("error update cache. %w", err)),
		)
	}

	return quote, nil
}

func (c *currencyInteractor) fetchQuote(
	ctx context.Context,
	cdrUrl string,
	params GetQuotationParams,
) (*entity.CurrencyQuotation, error) {
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

	if !data.Success {
		return nil, fmt.Errorf("eror fetch currency quote from api. %s", data.Description)
	}

	return &entity.CurrencyQuotation{
		Date:         params.Date,
		CurrencyID:   params.CurrencyID,
		CurrencyIDTo: params.CurrencyIDTo,
		Value:        data.Rates[params.CurrencyIDTo],
	}, nil
}

func (c *currencyInteractor) RunWorker() {
	go func() {
		c.log.Info("currency quotations collecting worker started")

		defer func() {
			c.log.Info("stop currency quotations collecting worker! bye bye!")

			if err := recover(); err != nil {
				c.log.Error("worker panic!", slog.Any("data", err))
			}
		}()

		startedAt := time.Now().UTC()
		day := startedAt.Day()

		if startedAt.Hour() >= 10 {
			day += 1
		}

		beginWork := time.Date(
			startedAt.Year(),
			startedAt.Month(),
			day,
			10,
			00,
			0,
			0,
			time.UTC,
		)

		initialTicker := time.NewTicker(beginWork.Sub(startedAt))

		<-initialTicker.C
		initialTicker.Stop()

		c.log.Info(
			"worker main loop started",
		)

		t := time.NewTicker(time.Hour * 24)

		for {
			for _, curr := range _currenciesList {
				date := time.Now()

				url, err := c.requestCurrencyQuotationsURL(
					date,
					"RUB",
					curr,
				)
				if err != nil {
					c.log.Error(
						"error build request url",
						logger.Err(err),
					)

					break
				}

				func() {
					ctx, cancel := context.WithTimeout(
						context.TODO(),
						time.Second*10,
					)

					defer cancel()

					quote, err := c.fetchQuote(ctx, url, GetQuotationParams{
						Date:         date,
						CurrencyID:   curr,
						CurrencyIDTo: "RUB",
					})
					if err != nil {
						c.log.Error(
							"error fetch quote",
							logger.Err(err),
						)

						return
					}

					c.log.Debug(
						"curr cached",
						slog.Any("quote", quote),
					)

					if err = c.cache.Set(
						ctx,
						quote.CurrencyID+quote.CurrencyIDTo+buildDateKey(date),
						*quote,
						time.Hour*24,
					); err != nil {
						c.log.Error(
							"error update cache",
							logger.Err(
								err,
							),
						)
					}

				}()
			}

			<-t.C
			t.Reset(time.Hour * 24)
		}
	}()
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
	q.Set("base", currFrom)

	if currTo != "" {
		q.Set("currencies", currTo)
	}

	return u.String() + "?" + q.Encode(), nil
}

func buildDateKey(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}

	layout := "2006-01-02"

	return t.Format(layout)
}
