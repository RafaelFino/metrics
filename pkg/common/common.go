package common

import "time"

type MetricType string

const (
	CounterType   MetricType = "counter"
	GaugeType     MetricType = "gauge"
	HistogramType MetricType = "histogram"
	SummaryType   MetricType = "summary"
)

type Metric struct {
	Name  string
	When  time.Time
	Value float64
	Tags  map[string]string
	Type  MetricType
}

type Series interface {
	Increment(value float64)
	Count() int
	Sum() float64
	Avg() float64
	Data() []float64
	CreatedAt() time.Time
	LastAt() time.Time
	Tags() map[string]string
	Name() string
}

type HistogramResolution string

const (
	SecondResolution HistogramResolution = "second"
	MinuteResolution HistogramResolution = "minute"
	HourResolution   HistogramResolution = "hour"
	DayResolution    HistogramResolution = "day"
)

type Exporter interface {
	Export(data *MetricData) error
}

type MetricData struct {
	Counters   map[string]*Metric
	Gauges     map[string]*Metric
	Histograms map[string]Series
	Summaries  map[string]Series
}
