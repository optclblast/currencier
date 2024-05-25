package repo

import (
	"github.com/optclblast/currencier/internal/pkg/postgres"
)

const defaultEntityCap = 64

// CurrencyRepo.
type CurrencyRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *CurrencyRepo {
	return &CurrencyRepo{pg}
}
