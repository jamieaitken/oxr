package oxr

import (
	"strings"
	"time"
)

type timeSeriesParams struct {
	startDate             time.Time
	endDate               time.Time
	baseCurrency          string
	destinationCurrencies string
	showAlternative       bool
	prettyPrint           bool
}

// TimeSeriesOption allows the client to specify values for a TimeSeries request.
type TimeSeriesOption func(params *timeSeriesParams)

// TimeSeriesForBaseCurrency sets the base currency to be used.
func TimeSeriesForBaseCurrency(currency string) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.baseCurrency = currency
	}
}

// TimeSeriesForDestinationCurrencies sets the destination currencies to be included in the response.
func TimeSeriesForDestinationCurrencies(currencies []string) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.destinationCurrencies = strings.Join(currencies, ",")
	}
}

// TimeSeriesWithAlternatives sets whether to include alternative currencies in the response.
func TimeSeriesWithAlternatives(active bool) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.showAlternative = active
	}
}

// TimeSeriesForStartDate sets the start date of the period.
func TimeSeriesForStartDate(start time.Time) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.startDate = start
	}
}

// TimeSeriesForEndDate sets the end date of the period.
func TimeSeriesForEndDate(end time.Time) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.endDate = end
	}
}

// TimeSeriesWithPrettyPrint sets whether to minify the response.
func TimeSeriesWithPrettyPrint(active bool) TimeSeriesOption {
	return func(p *timeSeriesParams) {
		p.prettyPrint = active
	}
}
