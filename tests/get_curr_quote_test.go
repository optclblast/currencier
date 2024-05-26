package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type apiResponse struct {
	Date         time.Time `json:"date"`
	CurrencyID   string    `json:"currency_id"`
	CurrencyIDTo string    `json:"currency_id_to"`
	Value        float64   `json:"value"`

	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"message"`
}

const defaultHost = "http://localhost:8080"

func TestGetUSDQuoteByDate(t *testing.T) {
	host := os.Getenv("CURRENCIER_TESTS_HOST")

	if host == "" {
		host = defaultHost
	}

	baseUrl := host + "/currency"

	tests := []struct {
		name         string
		url          string
		currencyID   string
		currencyIDTo string
		expectedVal  float64
		mustFail     bool
	}{
		// Valid cases, shoud pass
		{
			name:         "EUR to RUB. Correct url (25.05.2024)",
			url:          baseUrl + "?date=25.05.2024&val=EUR",
			currencyID:   "EUR",
			currencyIDTo: "RUB",
			expectedVal:  97.207957,
			mustFail:     false,
		},
		{
			name:         "USD to RUB. Correct url (25.05.2024)",
			url:          baseUrl + "?date=25.05.2024&val=USD",
			currencyID:   "USD",
			currencyIDTo: "RUB",
			expectedVal:  89.574226,
			mustFail:     false,
		},
		{
			name:         "USD to JPY. Correct url (25.05.2024)",
			url:          baseUrl + "?date=25.05.2024&val=USD&val_to=JPY",
			currencyID:   "USD",
			currencyIDTo: "JPY",
			expectedVal:  156.972537,
			mustFail:     false,
		},

		// Invalid cases, shoud fail
		{
			name:     "EUR to RUB. Invalid day (111.05.2023)",
			url:      baseUrl + "?date=111.05.2023&val=EUR",
			mustFail: true,
		},
		{
			name:     "EUR to RUB. Invalid month (01.13.2024)",
			url:      baseUrl + "?date=01.13.2024&val=EUR",
			mustFail: true,
		},
		{
			name:     "EUR to RUB. Invalid year. Year out of range (25.05.2025)",
			url:      baseUrl + "?date=25.05.2025&val=EUR",
			mustFail: true,
		},
		{
			name:     "EUR to RUB. Invalid year. Zero year (25.05.0000)",
			url:      baseUrl + "?date=25.05.2025&val=EUR",
			mustFail: true,
		},
		{
			name:     "USD to EUR. Invalid year. Zero year (25.05.0000)",
			url:      baseUrl + "?date=25.05.2025&val=EUR",
			mustFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.Get(test.url)
			if err != nil {
				t.Fatalf("error fetch data from api. %s", err.Error())
			}

			defer r.Body.Close()

			respRaw, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("error read api response. %s", err.Error())
			}

			resp := new(apiResponse)

			if err := json.Unmarshal(respRaw, resp); err != nil {
				t.Fatalf("error unmarshal api response. %s", err.Error())
			}

			if !test.mustFail && resp.ErrorCode > 0 || !test.mustFail && r.StatusCode != http.StatusOK {
				t.Fatalf("request faild, but it shouldn't")
			}

			if resp.Value != test.expectedVal {
				t.Fatalf("incorrect response currency value. expected: %v, got: %v", test.expectedVal, resp.Value)
			}

			if !test.mustFail {
				if test.currencyID != resp.CurrencyID {
					t.Fatalf("incorrect response currency id. expected: %s, got: %s", test.currencyID, resp.CurrencyID)
				}

				if test.currencyIDTo != resp.CurrencyIDTo {
					t.Fatalf("incorrect response currency id to. expected: %s, got: %s", test.currencyIDTo, resp.CurrencyIDTo)
				}
			}
		})
	}
}
