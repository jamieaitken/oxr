# OXR 💹

## Install

```shell
go get github.com/jamieaitken/oxr
```

## How to use

### Initialise your client

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))
```

### Latest Rates

Latest retrieves the latest exchange rates available from the Open Exchange
Rates [API](https://docs.openexchangerates.org/docs/latest-json)

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

latestRates, err := c.Latest(context.Background(), oxr.LatestForBaseCurrency("GBP"))
```

### Historical

Historical retrieves historical exchange rates for any date available from the Open Exchange Rates
[API](https://docs.openexchangerates.org/docs/historical-json), currently going back to 1st January 1999.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

historicalRates, err := c.Historical(
context.Background(),
oxr.HistoricalForDate(time.Date(2022, 03, 10, 12, 00, 00, 00, time.UTC)),
oxr.HistoricalForBaseCurrency("USD"),
)
```

### Currencies

Currencies retrieves the list of all currency symbols available from the Open Exchange
Rates [API](https://docs.openexchangerates.org/docs/currencies-json), along with their full names.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

currencies, err := c.Currencies(context.Background())
```

### Time Series

TimeSeries retrieves historical exchange rates for a given time period, where available, using the time series / bulk
download [API](https://docs.openexchangerates.org/docs/time-series-json) endpoint.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

currencies, err := c.TimeSeries(
context.Background(),
oxr.TimeSeriesForStartDate(time.Date(2013, 01, 01, 00, 00, 00, 00, time.UTC)),
oxr.TimeSeriesForEndDate(time.Date(2013, 01, 31, 00, 00, 00, 00, time.UTC)),
oxr.TimeSeriesForBaseCurrency("AUD"),
oxr.TimeSeriesForDestinationCurrencies([]string{"BTC", "EUR", "HKD"}),
)
```

### Convert

Convert any money value from one currency to another at the
latest [API](https://docs.openexchangerates.org/docs/convert) rates.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

currencies, err := c.Convert(
context.Background(),
oxr.ConvertWithValue(100.12),
oxr.ConvertForBaseCurrency("GBP"),
oxr.ConvertForDestinationCurrency("USD"),
)
```

### Open High Low Close (OHLC)

OpenHighLowClose [retrieves](https://docs.openexchangerates.org/docs/ohlc-json) historical Open, High Low, Close (OHLC)
and Average exchange rates for a given time period, ranging from 1 month to 1 minute, where available.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

currencies, err := c.OpenHighLowClose(
context.Background(),
oxr.OHLCForBaseCurrency("USD"),
oxr.OHLCForPeriod(oxr.ThirtyMinute),
oxr.OHLCForDestinationCurrencies([]string{"GBP", "EUR"}),
oxr.OHLCForStartTime(time.Date(2022, 3, 15, 13, 00, 00, 00, time.UTC)),
)
```

### Usage

Usage [retrieves](https://docs.openexchangerates.org/docs/usage-json) basic plan information and usage statistics for an
Open Exchange Rates App ID.

```go
doer := http.DefaultClient
c := oxr.New(oxr.WithAppID("your_app_id"), oxr.WithDoer(doer))

currencies, err := c.Usage(context.Background())
```

## License
[MIT](https://choosealicense.com/licenses/mit/)