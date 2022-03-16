package oxr

import (
	"strings"
	"time"
)

type historicalParams struct {
	date                  time.Time
	baseCurrency          string
	destinationCurrencies string
	showAlternative       bool
	prettyPrint           bool
}

// HistoricalOption allows the client to specify values for a historical request.
type HistoricalOption func(*historicalParams)

// HistoricalForBaseCurrency sets the base currency for the historical request.
func HistoricalForBaseCurrency(currency string) HistoricalOption {
	return func(p *historicalParams) {
		p.baseCurrency = currency
	}
}

// HistoricalForDestinationCurrencies sets the destination currency for the historical request.
func HistoricalForDestinationCurrencies(currencies []string) HistoricalOption {
	return func(p *historicalParams) {
		p.destinationCurrencies = strings.Join(currencies, ",")
	}
}

// HistoricalWithAlternatives sets whether to include alternative currencies.
func HistoricalWithAlternatives(active bool) HistoricalOption {
	return func(p *historicalParams) {
		p.showAlternative = active
	}
}

// HistoricalForDate sets the date of the rate to be requested.
func HistoricalForDate(date time.Time) HistoricalOption {
	return func(p *historicalParams) {
		p.date = date
	}
}

// HistoricalWithPrettyPrint sets whether to minify the response.
func HistoricalWithPrettyPrint(active bool) HistoricalOption {
	return func(p *historicalParams) {
		p.prettyPrint = active
	}
}
