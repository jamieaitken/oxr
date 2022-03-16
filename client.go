package oxr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02"
	basePath   = "https://openexchangerates.org/api/"
)

var (
	ErrBadResponse = errors.New("failed to receive successful response")
)

// Doer sends a http.Request and returns a http.Response.
type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client is responsible for all interactions between OXR.
type Client struct {
	appID   string
	doer    Doer
	baseURL string
}

// New instantiates a Client.
func New(opts ...ClientOption) Client {
	c := Client{
		baseURL: basePath,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// Latest retrieves the latest exchange rates available from the Open Exchange Rates API.
// https://docs.openexchangerates.org/docs/latest-json
func (c Client) Latest(ctx context.Context, opts ...LatestOption) (LatestRatesResponse, error) {
	r := latestParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%slatest.json", c.baseURL), http.NoBody)
	if err != nil {
		return LatestRatesResponse{}, err
	}

	v := req.URL.Query()

	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))
	if r.baseCurrency != "" {
		v.Add("base", r.baseCurrency)
	}
	if r.destinationCurrencies != "" {
		v.Add("symbols", r.destinationCurrencies)
	}
	v.Add("show_alternative", strconv.FormatBool(r.showAlternative))

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()
	if err != nil {
		return LatestRatesResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return LatestRatesResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData LatestRatesResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return LatestRatesResponse{}, err
	}

	return resData, nil
}

// Historical retrieves historical exchange rates for any date available from the Open Exchange Rates API, currently
// going back to 1st January 1999.
// https://docs.openexchangerates.org/docs/historical-json
func (c Client) Historical(ctx context.Context, opts ...HistoricalOption) (HistoricalRatesResponse, error) {
	r := historicalParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%shistorical/%s.json", c.baseURL, r.date.Format(timeFormat)), http.NoBody)
	if err != nil {
		return HistoricalRatesResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))
	v.Add("show_alternative", strconv.FormatBool(r.showAlternative))
	if r.baseCurrency != "" {
		v.Add("base", r.baseCurrency)
	}
	if r.destinationCurrencies != "" {
		v.Add("symbols", r.destinationCurrencies)
	}

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return HistoricalRatesResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return HistoricalRatesResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData HistoricalRatesResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return HistoricalRatesResponse{}, err
	}

	return resData, nil
}

// Currencies retrieves the list of all currency symbols available from the Open Exchange Rates API, along with their
// full names.
// https://docs.openexchangerates.org/docs/currencies-json
func (c Client) Currencies(ctx context.Context, opts ...CurrenciesOption) (CurrenciesResponse, error) {
	r := currenciesParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%scurrencies.json", c.baseURL), http.NoBody)
	if err != nil {
		return CurrenciesResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("show_inactive", strconv.FormatBool(r.showInactive))
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))
	v.Add("show_alternative", strconv.FormatBool(r.showAlternative))

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return CurrenciesResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return CurrenciesResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData CurrenciesResponse
	err = json.NewDecoder(res.Body).Decode(&resData.Currencies)
	if err != nil {
		return CurrenciesResponse{}, err
	}

	return resData, nil
}

// TimeSeries retrieves historical exchange rates for a given time period, where available, using the time series / bulk
// download API endpoint.
// https://docs.openexchangerates.org/docs/time-series-json
func (c Client) TimeSeries(ctx context.Context, opts ...TimeSeriesOption) (TimeSeriesResponse, error) {
	r := timeSeriesParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%stime-series.json", c.baseURL), http.NoBody)
	if err != nil {
		return TimeSeriesResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))
	v.Add("show_alternative", strconv.FormatBool(r.showAlternative))
	v.Add("start", r.startDate.Format(timeFormat))
	v.Add("end", r.endDate.Format(timeFormat))
	if r.destinationCurrencies != "" {
		v.Add("symbols", r.destinationCurrencies)
	}
	if r.baseCurrency != "" {
		v.Add("base", r.baseCurrency)
	}

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return TimeSeriesResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return TimeSeriesResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData TimeSeriesResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return TimeSeriesResponse{}, err
	}

	return resData, nil
}

// Convert any money value from one currency to another at the latest API rates.
// https://docs.openexchangerates.org/docs/convert
func (c Client) Convert(ctx context.Context, opts ...ConvertOption) (ConversionResponse, error) {
	r := convertParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%sconvert/%v/%s/%s", c.baseURL, r.value, r.baseCurrency, r.destinationCurrency), http.NoBody)
	if err != nil {
		return ConversionResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return ConversionResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return ConversionResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData ConversionResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return ConversionResponse{}, err
	}

	return resData, nil
}

// OpenHighLowClose retrieves historical Open, High Low, Close (OHLC) and Average exchange rates for a given time period,
// ranging from 1 month to 1 minute, where available.
// https://docs.openexchangerates.org/docs/ohlc-json
func (c Client) OpenHighLowClose(ctx context.Context, opts ...OHLCOption) (OHLCResponse, error) {
	r := ohlcParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%sohlc.json", c.baseURL), http.NoBody)
	if err != nil {
		return OHLCResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))
	v.Add("start_date", r.startTime.Format(time.RFC3339))
	v.Add("period", r.period.String())
	if r.destinationCurrencies != "" {
		v.Add("symbols", r.destinationCurrencies)
	}
	if r.baseCurrency != "" {
		v.Add("base", r.baseCurrency)
	}

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return OHLCResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return OHLCResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData OHLCResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return OHLCResponse{}, err
	}

	return resData, nil
}

// Usage retrieves basic plan information and usage statistics for an Open Exchange Rates App ID.
// https://docs.openexchangerates.org/docs/usage-json
func (c Client) Usage(ctx context.Context, opts ...UsageOption) (UsageResponse, error) {
	r := usageParams{}

	for _, opt := range opts {
		opt(&r)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%susage.json", c.baseURL), http.NoBody)
	if err != nil {
		return UsageResponse{}, err
	}

	v := req.URL.Query()
	v.Add("app_id", c.appID)
	v.Add("prettyprint", strconv.FormatBool(r.prettyPrint))

	req.URL.RawQuery = v.Encode()

	res, err := c.doer.Do(req)
	defer res.Body.Close()

	if err != nil {
		return UsageResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return UsageResponse{}, fmt.Errorf("status received: %v: %w", res.StatusCode, ErrBadResponse)
	}

	var resData UsageResponse
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		return UsageResponse{}, err
	}

	return resData, nil
}
