package oxr

import "time"

// LatestRatesResponse is the response of a Latest request.
type LatestRatesResponse struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// ConversionResponse is the response of a Conversion request.
type ConversionResponse struct {
	Disclaimer string            `json:"disclaimer"`
	License    string            `json:"license"`
	Request    ConversionRequest `json:"request"`
	Meta       ConversionMeta    `json:"meta"`
	Response   float64           `json:"response"`
}

type ConversionRequest struct {
	Query  string  `json:"query"`
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}

type ConversionMeta struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}

// CurrenciesResponse is the response of a Currencies request.
type CurrenciesResponse struct {
	Currencies map[string]string
}

// HistoricalRatesResponse is the response of a Historical request.
type HistoricalRatesResponse struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// OHLCResponse is the response of a OHLC request.
type OHLCResponse struct {
	Disclaimer string              `json:"disclaimer"`
	License    string              `json:"license"`
	StartTime  time.Time           `json:"start_time"`
	EndTime    time.Time           `json:"end_time"`
	Base       string              `json:"base"`
	Rates      map[string]OHLCRate `json:"rates"`
}

type OHLCRate struct {
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
	Average float64 `json:"average"`
}

// TimeSeriesResponse is the response of a TimeSeries request.
type TimeSeriesResponse struct {
	Disclaimer string                        `json:"disclaimer"`
	License    string                        `json:"license"`
	StartDate  string                        `json:"start_date"`
	EndDate    string                        `json:"end_date"`
	Base       string                        `json:"base"`
	Rates      map[string]map[string]float64 `json:"rates"`
}

// UsageResponse is the response of a Usage request.
type UsageResponse struct {
	Status int       `json:"status"`
	Data   UsageData `json:"data"`
}

type UsageData struct {
	AppID  string        `json:"app_id"`
	Status string        `json:"status"`
	Plan   UsageDataPlan `json:"plan"`
	Usage  DataUsage     `json:"usage"`
}

type DataUsage struct {
	Requests          int `json:"requests"`
	RequestsQuota     int `json:"requests_quota"`
	RequestsRemaining int `json:"requests_remaining"`
	DaysElapsed       int `json:"days_elapsed"`
	DaysRemaining     int `json:"days_remaining"`
	DailyAverage      int `json:"daily_average"`
}

type UsageDataPlan struct {
	Name            string                `json:"name"`
	Quota           string                `json:"quota"`
	UpdateFrequency string                `json:"update_frequency"`
	Features        UsageDataPlanFeatures `json:"features"`
}

type UsageDataPlanFeatures struct {
	Base         bool `json:"base"`
	Symbols      bool `json:"symbols"`
	Experimental bool `json:"experimental"`
	TimeSeries   bool `json:"time-series"`
	Convert      bool `json:"convert"`
}
