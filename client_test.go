package oxr_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jamieaitken/oxr"
)

func TestClient_Convert_Success(t *testing.T) {
	tests := []struct {
		name             string
		givenDoer        *mockDoer
		givenClientOpts  []oxr.ClientOption
		givenConvertOpts []oxr.ConvertOption
		expectedURL      string
		expectedResult   oxr.ConversionResponse
	}{
		{
			name: "given successful convert response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulConversion())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenConvertOpts: []oxr.ConvertOption{
				oxr.ConvertWithValue(100.12),
				oxr.ConvertForBaseCurrency("GBP"),
				oxr.ConvertForDestinationCurrency("USD"),
			},
			expectedURL: "https://openexchangerates.org/api/convert/100.12/GBP/USD?app_id=test&prettyprint=false",
			expectedResult: oxr.ConversionResponse{
				Disclaimer: "https://openexchangerates.org/terms/",
				License:    "https://openexchangerates.org/license/",
				Request: oxr.ConversionRequest{
					Query:  "/convert/100.12/GBP/USD",
					Amount: 100.12,
					From:   "GBP",
					To:     "USD",
				},
				Meta: oxr.ConversionMeta{
					Timestamp: 1449885661,
					Rate:      0.76,
				},
				Response: 76.0912,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.Convert(context.Background(), test.givenConvertOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_Convert_Fail(t *testing.T) {
	tests := []struct {
		name             string
		givenDoer        *mockDoer
		givenClientOpts  []oxr.ClientOption
		givenConvertOpts []oxr.ConvertOption
		expectedURL      string
		expectedError    error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenConvertOpts: []oxr.ConvertOption{
				oxr.ConvertWithValue(100.12),
				oxr.ConvertForBaseCurrency("GBP"),
				oxr.ConvertForDestinationCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/convert/100.12/GBP/USD?app_id=test&prettyprint=false",
			expectedError: http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenConvertOpts: []oxr.ConvertOption{
				oxr.ConvertWithValue(100.12),
				oxr.ConvertForBaseCurrency("GBP"),
				oxr.ConvertForDestinationCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/convert/100.12/GBP/USD?app_id=test&prettyprint=false",
			expectedError: oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.Convert(context.Background(), test.givenConvertOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_Currencies_Success(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenCurrenciesOpts []oxr.CurrenciesOption
		expectedURL         string
		expectedResult      oxr.CurrenciesResponse
	}{
		{
			name: "given successful currencies response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulCurrencies())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenCurrenciesOpts: []oxr.CurrenciesOption{
				oxr.CurrenciesWithInactive(true),
				oxr.CurrenciesWithPrettyPrint(true),
			},
			expectedURL: "https://openexchangerates.org/api/currencies.json?app_id=test&prettyprint=true&show_alternative=false&show_inactive=true",
			expectedResult: oxr.CurrenciesResponse{
				Currencies: map[string]string{
					"EUR": "Euro",
					"GBP": "Pound sterling",
					"USD": "US Dollar",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.Currencies(context.Background(), test.givenCurrenciesOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_Currencies_Fail(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenCurrenciesOpts []oxr.CurrenciesOption
		expectedURL         string
		expectedError       error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenCurrenciesOpts: []oxr.CurrenciesOption{},
			expectedURL:         "https://openexchangerates.org/api/currencies.json?app_id=test&prettyprint=false&show_alternative=false&show_inactive=false",
			expectedError:       http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenCurrenciesOpts: []oxr.CurrenciesOption{},
			expectedURL:         "https://openexchangerates.org/api/currencies.json?app_id=test&prettyprint=false&show_alternative=false&show_inactive=false",
			expectedError:       oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.Currencies(context.Background(), test.givenCurrenciesOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_Historical_Success(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenHistoricalOpts []oxr.HistoricalOption
		expectedURL         string
		expectedResult      oxr.HistoricalRatesResponse
	}{
		{
			name: "given successful historical response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulHistorical())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenHistoricalOpts: []oxr.HistoricalOption{
				oxr.HistoricalForDate(time.Date(2022, 3, 10, 12, 0, 0, 0, time.UTC)),
				oxr.HistoricalForBaseCurrency("USD"),
			},
			expectedURL: "https://openexchangerates.org/api/historical/2022-03-10.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedResult: oxr.HistoricalRatesResponse{
				Disclaimer: "Usage subject to terms: https://openexchangerates.org/terms",
				License:    "https://openexchangerates.org/license",
				Timestamp:  1341936000,
				Base:       "USD",
				Rates: map[string]float64{
					"GBP": 0.76,
					"EUR": 0.93,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.Historical(context.Background(), test.givenHistoricalOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_Historical_Fail(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenHistoricalOpts []oxr.HistoricalOption
		expectedURL         string
		expectedError       error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenHistoricalOpts: []oxr.HistoricalOption{
				oxr.HistoricalForDate(time.Date(2022, 3, 10, 12, 0, 0, 0, time.UTC)),
				oxr.HistoricalForBaseCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/historical/2022-03-10.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedError: http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenHistoricalOpts: []oxr.HistoricalOption{
				oxr.HistoricalForDate(time.Date(2022, 3, 10, 12, 0, 0, 0, time.UTC)),
				oxr.HistoricalForBaseCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/historical/2022-03-10.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedError: oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.Historical(context.Background(), test.givenHistoricalOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_Latest_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenLatestOpts []oxr.LatestOption
		expectedURL     string
		expectedResult  oxr.LatestRatesResponse
	}{
		{
			name: "given successful latest response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulLatest())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenLatestOpts: []oxr.LatestOption{
				oxr.LatestForBaseCurrency("USD"),
			},
			expectedURL: "https://openexchangerates.org/api/latest.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedResult: oxr.LatestRatesResponse{
				Disclaimer: "Usage subject to terms: https://openexchangerates.org/terms",
				License:    "https://openexchangerates.org/license",
				Timestamp:  1647453600,
				Base:       "USD",
				Rates: map[string]float64{
					"GBP": 0.764018,
					"KRW": 1225.826828,
					"USD": 1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.Latest(context.Background(), test.givenLatestOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_Latest_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenLatestOpts []oxr.LatestOption
		expectedURL     string
		expectedError   error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenLatestOpts: []oxr.LatestOption{
				oxr.LatestForBaseCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/latest.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedError: http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenLatestOpts: []oxr.LatestOption{
				oxr.LatestForBaseCurrency("USD"),
			},
			expectedURL:   "https://openexchangerates.org/api/latest.json?app_id=test&base=USD&prettyprint=false&show_alternative=false",
			expectedError: oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.Latest(context.Background(), test.givenLatestOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_OpenHighLowClose_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenOHLCOpts   []oxr.OHLCOption
		expectedURL     string
		expectedResult  oxr.OHLCResponse
	}{
		{
			name: "given successful ohlc response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulOHLC())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenOHLCOpts: []oxr.OHLCOption{
				oxr.OHLCForBaseCurrency("USD"),
				oxr.OHLCForPeriod(oxr.ThirtyMinute),
				oxr.OHLCForDestinationCurrencies([]string{"GBP", "EUR"}),
				oxr.OHLCForStartTime(time.Date(2022, 3, 15, 13, 0, 0, 0, time.UTC)),
			},
			expectedURL: "https://openexchangerates.org/api/ohlc.json?app_id=test&base=USD&period=30m&prettyprint=false&start_date=2022-03-15T13%3A00%3A00Z&symbols=GBP%2CEUR",
			expectedResult: oxr.OHLCResponse{
				Disclaimer: "Usage subject to terms: https://openexchangerates.org/terms",
				License:    "https://openexchangerates.org/license",
				Base:       "USD",
				StartTime:  time.Date(2022, 3, 15, 13, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2022, 3, 15, 13, 30, 0, 0, time.UTC),
				Rates: map[string]oxr.OHLCRate{
					"EUR": {
						Open:    0.872674,
						High:    0.872674,
						Low:     0.87203,
						Close:   0.872251,
						Average: 0.872253,
					},
					"GBP": {
						Open:    0.765284,
						High:    0.7657,
						Low:     0.7652,
						Close:   0.765541,
						Average: 0.765503,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.OpenHighLowClose(context.Background(), test.givenOHLCOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_OpenHighLowClose_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenOHLCOpts   []oxr.OHLCOption
		expectedURL     string
		expectedError   error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenOHLCOpts: []oxr.OHLCOption{},
			expectedURL:   "https://openexchangerates.org/api/ohlc.json?app_id=test&period=&prettyprint=false&start_date=0001-01-01T00%3A00%3A00Z",
			expectedError: http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenOHLCOpts: []oxr.OHLCOption{},
			expectedURL:   "https://openexchangerates.org/api/ohlc.json?app_id=test&period=&prettyprint=false&start_date=0001-01-01T00%3A00%3A00Z",
			expectedError: oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.OpenHighLowClose(context.Background(), test.givenOHLCOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_TimeSeries_Success(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenTimeSeriesOpts []oxr.TimeSeriesOption
		expectedURL         string
		expectedResult      oxr.TimeSeriesResponse
	}{
		{
			name: "given successful time series response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulTimeSeries())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenTimeSeriesOpts: []oxr.TimeSeriesOption{
				oxr.TimeSeriesForStartDate(time.Date(2013, 1, 1, 13, 0, 0, 0, time.UTC)),
				oxr.TimeSeriesForEndDate(time.Date(2013, 1, 31, 13, 0, 0, 0, time.UTC)),
				oxr.TimeSeriesForBaseCurrency("AUD"),
				oxr.TimeSeriesForDestinationCurrencies([]string{"BTC", "EUR", "HKD"}),
				oxr.TimeSeriesWithPrettyPrint(true),
			},
			expectedURL: "https://openexchangerates.org/api/time-series.json?app_id=test&base=AUD&end=2013-01-31&prettyprint=true&show_alternative=false&start=2013-01-01&symbols=BTC%2CEUR%2CHKD",
			expectedResult: oxr.TimeSeriesResponse{
				Disclaimer: "Usage subject to terms: https://openexchangerates.org/terms/",
				License:    "https://openexchangerates.org/license/",
				Base:       "AUD",
				StartDate:  "2013-01-01",
				EndDate:    "2013-01-31",
				Rates: map[string]map[string]float64{
					"2013-01-01": {
						"BTC": 0.0778595876,
						"EUR": 0.785518,
						"HKD": 8.04136,
					},
					"2013-01-02": {
						"BTC": 0.0789400739,
						"EUR": 0.795034,
						"HKD": 8.138096,
					},
					"2013-01-03": {
						"BTC": 0.0785299961,
						"EUR": 0.80092,
						"HKD": 8.116954,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.TimeSeries(context.Background(), test.givenTimeSeriesOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_TimeSeries_Fail(t *testing.T) {
	tests := []struct {
		name                string
		givenDoer           *mockDoer
		givenClientOpts     []oxr.ClientOption
		givenTimeSeriesOpts []oxr.TimeSeriesOption
		expectedURL         string
		expectedError       error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenTimeSeriesOpts: []oxr.TimeSeriesOption{},
			expectedURL:         "https://openexchangerates.org/api/time-series.json?app_id=test&end=0001-01-01&prettyprint=false&show_alternative=false&start=0001-01-01",
			expectedError:       http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenTimeSeriesOpts: []oxr.TimeSeriesOption{},
			expectedURL:         "https://openexchangerates.org/api/time-series.json?app_id=test&end=0001-01-01&prettyprint=false&show_alternative=false&start=0001-01-01",
			expectedError:       oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.TimeSeries(context.Background(), test.givenTimeSeriesOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestClient_Usage_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenUsageOpts  []oxr.UsageOption
		expectedURL     string
		expectedResult  oxr.UsageResponse
	}{
		{
			name: "given successful usage response, expect payload returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(successfulUsage())),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenUsageOpts: []oxr.UsageOption{
				oxr.UsageWithPrettyPrint(true),
			},
			expectedURL: "https://openexchangerates.org/api/usage.json?app_id=test&prettyprint=true",
			expectedResult: oxr.UsageResponse{
				Status: 200,
				Data: oxr.UsageData{
					AppID:  "YOUR_APP_ID",
					Status: "active",
					Plan: oxr.UsageDataPlan{
						Name:            "Enterprise",
						Quota:           "100,000 requests/month",
						UpdateFrequency: "30-minute",
						Features: oxr.UsageDataPlanFeatures{
							Base:         true,
							Symbols:      true,
							Experimental: true,
							TimeSeries:   true,
							Convert:      false,
						},
					},
					Usage: oxr.DataUsage{
						Requests:          54524,
						RequestsQuota:     100000,
						RequestsRemaining: 45476,
						DaysElapsed:       16,
						DaysRemaining:     14,
						DailyAverage:      2842,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			actual, err := c.Usage(context.Background(), test.givenUsageOpts...)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(actual, test.expectedResult) {
				t.Fatal(cmp.Diff(actual, test.expectedResult))
			}
		})
	}
}

func TestClient_Usage_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenDoer       *mockDoer
		givenClientOpts []oxr.ClientOption
		givenUsageOpts  []oxr.UsageOption
		expectedURL     string
		expectedError   error
	}{
		{
			name: "given doer error, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
				GivenError: http.ErrBodyNotAllowed,
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenUsageOpts: []oxr.UsageOption{},
			expectedURL:    "https://openexchangerates.org/api/usage.json?app_id=test&prettyprint=false",
			expectedError:  http.ErrBodyNotAllowed,
		},
		{
			name: "given non 200 response, expect error returned",
			givenDoer: &mockDoer{
				GivenResponse: &http.Response{
					Status:     http.StatusText(http.StatusForbidden),
					StatusCode: http.StatusForbidden,
					Body:       io.NopCloser(strings.NewReader("")),
				},
			},
			givenClientOpts: []oxr.ClientOption{
				oxr.WithAppID("test"),
			},
			givenUsageOpts: []oxr.UsageOption{},
			expectedURL:    "https://openexchangerates.org/api/usage.json?app_id=test&prettyprint=false",
			expectedError:  oxr.ErrBadResponse,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := oxr.New(append(test.givenClientOpts, oxr.WithDoer(test.givenDoer))...)

			_, err := c.Usage(context.Background(), test.givenUsageOpts...)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(test.givenDoer.SpyURL, test.expectedURL) {
				t.Fatal(cmp.Diff(test.givenDoer.SpyURL, test.expectedURL))
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockDoer struct {
	GivenResponse *http.Response
	GivenError    error
	SpyURL        string
}

func (m *mockDoer) Do(r *http.Request) (*http.Response, error) {
	m.SpyURL = r.URL.String()

	return m.GivenResponse, m.GivenError
}

func successfulConversion() string {
	return `{
    "disclaimer": "https://openexchangerates.org/terms/",
    "license": "https://openexchangerates.org/license/",
    "request": {
        "query": "/convert/100.12/GBP/USD",
        "amount": 100.12,
        "from": "GBP",
        "to": "USD"
    },
    "meta": {
        "timestamp": 1449885661,
        "rate": 0.76
    },
    "response": 76.0912
}`
}

func successfulCurrencies() string {
	return `{
  "EUR": "Euro",
  "GBP": "Pound sterling",
  "USD": "US Dollar"
}`
}

func successfulHistorical() string {
	return `{
    "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
    "license": "https://openexchangerates.org/license",
    "timestamp": 1341936000,
    "base": "USD",
    "rates": {
        "GBP": 0.76,
		"EUR": 0.93
    }
}`
}

func successfulLatest() string {
	return `{
  "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
  "license": "https://openexchangerates.org/license",
  "timestamp": 1647453600,
  "base": "USD",
  "rates": {
    "GBP": 0.764018,
    "KRW": 1225.826828,
    "USD": 1
  }
}`
}

func successfulOHLC() string {
	return `{
  "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
  "license": "https://openexchangerates.org/license",
  "start_time": "2022-03-15T13:00:00Z",
  "end_time": "2022-03-15T13:30:00Z",
  "base": "USD",
  "rates": {
    "EUR": {
      "open": 0.872674,
      "high": 0.872674,
      "low": 0.87203,
      "close": 0.872251,
      "average": 0.872253
    },
    "GBP": {
      "open": 0.765284,
      "high": 0.7657,
      "low": 0.7652,
      "close": 0.765541,
      "average": 0.765503
    }
  }
}`
}

func successfulTimeSeries() string {
	return `{
    "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms/",
    "license": "https://openexchangerates.org/license/",
    "start_date": "2013-01-01",
    "end_date": "2013-01-31",
    "base": "AUD",
    "rates": {
        "2013-01-01": {
            "BTC": 0.0778595876,
            "EUR": 0.785518,
            "HKD": 8.04136
        },
        "2013-01-02": {
            "BTC": 0.0789400739,
            "EUR": 0.795034,
            "HKD": 8.138096
        },
        "2013-01-03": {
            "BTC": 0.0785299961,
            "EUR": 0.80092,
            "HKD": 8.116954
        }
    }
}`
}

func successfulUsage() string {
	return `{
  "status": 200,
  "data": {
    "app_id": "YOUR_APP_ID",
    "status": "active",
    "plan": {
      "name": "Enterprise",
      "quota": "100,000 requests/month",
      "update_frequency": "30-minute",
      "features": {
        "base": true,
        "symbols": true,
        "experimental": true,
        "time-series": true,
        "convert": false
      }
    },
    "usage": {
      "requests": 54524,
      "requests_quota": 100000,
      "requests_remaining": 45476,
      "days_elapsed": 16,
      "days_remaining": 14,
      "daily_average": 2842
    }
  }
}`
}
