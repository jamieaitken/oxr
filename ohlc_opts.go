package oxr

import (
	"strings"
	"time"
)

// Available periods.
const (
	OneMinute     period = "1m"
	FiveMinute    period = "5m"
	FifteenMinute period = "15m"
	ThirtyMinute  period = "30m"
	OneHour       period = "1h"
	TwelveHour    period = "12h"
	OneDay        period = "1d"
	OneWeek       period = "1w"
	OneMonth      period = "1mo"
)

type ohlcParams struct {
	startTime             time.Time
	period                period
	baseCurrency          string
	destinationCurrencies string
	prettyPrint           bool
}

type period string

// String implements a fmt.Stringer for period.
func (p period) String() string {
	return string(p)
}

// OHLCOption allows the client to specify values for a OHLC request.
type OHLCOption func(params *ohlcParams)

// OHLCWithPrettyPrint sets whether to minify the response.
func OHLCWithPrettyPrint(active bool) OHLCOption {
	return func(p *ohlcParams) {
		p.prettyPrint = active
	}
}

// OHLCForStartTime sets the start time for the given period.
func OHLCForStartTime(startTime time.Time) OHLCOption {
	return func(p *ohlcParams) {
		p.startTime = startTime
	}
}

// OHLCForPeriod sets the length of the period.
func OHLCForPeriod(period period) OHLCOption {
	return func(p *ohlcParams) {
		p.period = period
	}
}

// OHLCForBaseCurrency sets the base currency.
func OHLCForBaseCurrency(currency string) OHLCOption {
	return func(p *ohlcParams) {
		p.baseCurrency = currency
	}
}

// OHLCForDestinationCurrencies sets the destination currencies to be included in the response.
func OHLCForDestinationCurrencies(destinationCurrencies []string) OHLCOption {
	return func(p *ohlcParams) {
		p.destinationCurrencies = strings.Join(destinationCurrencies, ",")
	}
}
