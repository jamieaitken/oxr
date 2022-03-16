package oxr

import (
	"strings"
)

type latestParams struct {
	baseCurrency          string
	destinationCurrencies string
	showAlternative       bool
	prettyPrint           bool
}

// LatestOption allows the client to specify values for a latest request.
type LatestOption func(params *latestParams)

// LatestForBaseCurrency sets the base currency.
func LatestForBaseCurrency(currency string) LatestOption {
	return func(p *latestParams) {
		p.baseCurrency = currency
	}
}

// LatestForDestinationCurrencies sets the destination currencies to be included in the response.
func LatestForDestinationCurrencies(currencies []string) LatestOption {
	return func(p *latestParams) {
		p.destinationCurrencies = strings.Join(currencies, ",")
	}
}

// LatestWithAlternatives sets whether to include alternative currencies.
func LatestWithAlternatives(active bool) LatestOption {
	return func(p *latestParams) {
		p.showAlternative = active
	}
}

// LatestWithPrettyPrint sets whether to minify the response.
func LatestWithPrettyPrint(active bool) LatestOption {
	return func(p *latestParams) {
		p.prettyPrint = active
	}
}
