package oxr

type convertParams struct {
	value               float64
	baseCurrency        string
	destinationCurrency string
	prettyPrint         bool
}

// ConvertOption allows the client to specify values for a conversion request.
type ConvertOption func(*convertParams)

// ConvertForBaseCurrency sets the base currency for a conversion.
func ConvertForBaseCurrency(currency string) ConvertOption {
	return func(p *convertParams) {
		p.baseCurrency = currency
	}
}

// ConvertForDestinationCurrency sets the destination currency for a conversion.
func ConvertForDestinationCurrency(currency string) ConvertOption {
	return func(p *convertParams) {
		p.destinationCurrency = currency
	}
}

// ConvertWithValue sets the value to be converted.
func ConvertWithValue(value float64) ConvertOption {
	return func(p *convertParams) {
		p.value = value
	}
}

// ConvertWithPrettyPrint sets whether to minify the response.
func ConvertWithPrettyPrint(active bool) ConvertOption {
	return func(p *convertParams) {
		p.prettyPrint = active
	}
}
