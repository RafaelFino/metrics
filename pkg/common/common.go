package common

import "time"

type MetricType string

const (
	CounterType   MetricType = "counter"
	GaugeType     MetricType = "gauge"
	HistogramType MetricType = "histogram"
)

type Metric struct {
	Name  string
	When  time.Time
	Value float64
	Tags  map[string]string
	Type  MetricType
}

type WindowResolution string

const (
	SecondResolution WindowResolution = "second"
	MinuteResolution WindowResolution = "minute"
	HourResolution   WindowResolution = "hour"
	DayResolution    WindowResolution = "day"
)

type Exporter interface {
	Export(data *MetricData) error
}

type MetricData struct {
	Counters   map[string]*Metric
	Gauges     map[string]*Metric
	Histograms map[string]*Serie
}

type Serie struct {
	Name       string
	Type       MetricType
	Count      int
	Sum        float64
	Avg        float64
	CreatedAt  time.Time
	LastAt     time.Time
	Tags       map[string]string
	Data       map[int64]float64
	Resolution WindowResolution
}

func NewSerie(name string, metricType MetricType, tags map[string]string, resolution WindowResolution, value float64) *Serie {
	return &Serie{
		Name:       name,
		CreatedAt:  time.Now(),
		LastAt:     time.Now(),
		Tags:       tags,
		Data:       map[int64]float64{getKey(resolution): value},
		Resolution: resolution,
		Type:       metricType,
	}
}

func getKey(resolution WindowResolution) int64 {
	when := time.Now()

	switch resolution {
	case SecondResolution:
		when = when.Truncate(time.Second)
	case MinuteResolution:
		when = when.Truncate(time.Minute)
	case HourResolution:
		when = when.Truncate(time.Hour)
	case DayResolution:
		when = when.Truncate(time.Hour * 24)
	}

	return when.Unix()
}

func (s *Serie) Clear() {
	s.Data = make(map[int64]float64)
	s.CreatedAt = time.Now()
	s.LastAt = time.Now()
}

func (s *Serie) Increment(value float64) {
	s.Data[getKey(s.Resolution)] += value
	s.LastAt = time.Now()
}

func (s *Serie) Calculate() {
	s.Count = len(s.Data)
	s.Sum = 0

	for _, i := range s.Data {
		s.Sum += i
	}

	if s.Count != 0 {
		s.Avg = s.Sum / float64(s.Count)
	}
}
