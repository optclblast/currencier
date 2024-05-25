package entity

import (
	"encoding/json"
	"time"
)

type CurrencyQuotation struct {
	Date         time.Time `json:"date"`
	CurrencyID   string    `json:"currency_id"`
	CurrencyIDTo string    `json:"currency_id_to"`
	Value        float64   `json:"value"`
}

func (c CurrencyQuotation) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c CurrencyQuotation) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &c)
}
