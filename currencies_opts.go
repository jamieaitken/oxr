package oxr

type currenciesParams struct {
	showAlternative bool
	showInactive    bool
	prettyPrint     bool
}

// CurrenciesOption allows the client to specify values for a currencies request.
type CurrenciesOption func(params *currenciesParams)

// CurrenciesWithAlternatives includes alternative currencies.
func CurrenciesWithAlternatives(active bool) CurrenciesOption {
	return func(p *currenciesParams) {
		p.showAlternative = active
	}
}

// CurrenciesWithInactive includes historical/inactive currencies.
func CurrenciesWithInactive(active bool) CurrenciesOption {
	return func(p *currenciesParams) {
		p.showInactive = active
	}
}

// CurrenciesWithPrettyPrint sets whether to minify the response.
func CurrenciesWithPrettyPrint(active bool) CurrenciesOption {
	return func(p *currenciesParams) {
		p.prettyPrint = active
	}
}
