package repo

import (
	"context"
	"time"

	"github.com/optclblast/currencier/internal/entity"
	"github.com/optclblast/currencier/internal/pkg/postgres"
)

type GetParams struct {
	Date time.Time
	Val  string
}

type QutesRepo interface {
	Get(ctx context.Context, params GetParams) (float64, error)
	Set(ctx context.Context, quotation entity.CurrencyQuotation) error
}

// CurrencyRepo.
type qutesRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *qutesRepo {
	return &qutesRepo{pg}
}
