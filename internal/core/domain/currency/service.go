package currency

import "context"

//go:generate mockgen -package currencymocks -destination currencymocks/currency_service_mocks.go . Service

type Service interface {
	GetCurrencyPriceInDollar(ctx context.Context, currencyID string) (float64, error)
	IsValidCurrency(ctx context.Context, currencyID string) (bool, error)
}
