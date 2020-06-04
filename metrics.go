package metrics

import (
	"log"
	"time"
)

type MetricType string

const (
	CounterType   MetricType = "counter"
	GaugeType     MetricType = "gauge"
	HistogramType MetricType = "histogram"
	SummaryType   MetricType = "summary"
)

type Metric struct {
	Name      string
	CreatedAt int64
	LastAt    int64
	Value     float64
	Max       float64
	Min       float64
	Tags      map[string]string
	Type      MetricType
}

type WindowResolution string

const (
	SecondResolution WindowResolution = "second"
	MinuteResolution WindowResolution = "minute"
	HourResolution   WindowResolution = "hour"
	DayResolution    WindowResolution = "day"
)

type MetricData struct {
	Metrics map[string]*Metric
	Series  map[string]*Serie
}

type Serie struct {
	Name       string
	Type       MetricType
	Max        float64
	Min        float64
	Count      int
	Sum        float64
	Avg        float64
	CreatedAt  int64
	LastAt     int64
	Tags       map[string]string
	Data       map[int64]float64
	Resolution WindowResolution
}

type ExporterFunc = func(data *MetricData) error

func NewMetric(name string, metricType MetricType, tags map[string]string, value float64) *Metric {
	return &Metric{
		Name:      name,
		CreatedAt: time.Now().Unix(),
		LastAt:    time.Now().Unix(),
		Tags:      tags,
		Type:      metricType,
		Max:       value,
		Min:       value,
		Value:     value,
	}
}

func (m *Metric) Increment(v float64) {
	switch m.Type {
	case CounterType:
		m.Value += v
	case GaugeType:
		m.Value = v
	}

	if m.Max < v {
		m.Max = v
	}

	if m.Min > v {
		m.Min = v
	}

	m.LastAt = time.Now().Unix()
}

func NewSerie(name string, metricType MetricType, tags map[string]string, resolution WindowResolution, value float64) *Serie {
	return &Serie{
		Name:       name,
		CreatedAt:  time.Now().Unix(),
		LastAt:     time.Now().Unix(),
		Tags:       tags,
		Data:       map[int64]float64{getKey(resolution): value},
		Resolution: resolution,
		Type:       metricType,
		Max:        value,
		Min:        value,
	}
}

func getKey(resolution WindowResolution) int64 {
	when := time.Now()

	switch resolution {
	case SecondResolution:
		when = when.Round(time.Second)
	case MinuteResolution:
		when = when.Round(time.Minute)
	case HourResolution:
		when = when.Round(time.Hour)
	case DayResolution:
		when = when.Round(time.Hour * 24)
	}

	return when.Unix()
}

func (s *Serie) Clear() {
	s.Data = make(map[int64]float64)
	s.CreatedAt = time.Now().Unix()
	s.LastAt = time.Now().Unix()
}

func (s *Serie) Increment(value float64) {
	switch s.Type {
	case HistogramType:
		s.Data[getKey(s.Resolution)] = value
	case SummaryType:
		s.Data[getKey(s.Resolution)] += value
	}

	s.Count++
	s.LastAt = time.Now().Unix()
}

func (s *Serie) Calculate() {
	s.Sum = 0

	for _, i := range s.Data {
		if s.Max < i {
			s.Max = i
		}

		if s.Min > i {
			s.Min = i
		}

		s.Sum += i
	}

	if s.Count != 0 {
		s.Avg = s.Sum / float64(len(s.Data))
	}
}

type Processor struct {
	received    chan *Metric
	exportChan  chan *MetricData
	stopRequest chan bool

	metrics map[string]*Metric
	series  map[string]*Serie

	exporter       ExporterFunc
	exportInterval time.Duration

	lastExport *MetricData
}

func NewMetricProcessor(exportInterval time.Duration, exporter ExporterFunc) *Processor {
	ret := &Processor{
		received:       make(chan *Metric, 256),
		exportChan:     make(chan *MetricData, 10),
		stopRequest:    make(chan bool),
		exporter:       exporter,
		exportInterval: exportInterval,
		metrics:        make(map[string]*Metric),
		series:         make(map[string]*Serie),
		lastExport: &MetricData{
			Metrics: make(map[string]*Metric),
			Series:  make(map[string]*Serie),
		},
	}

	go ret.process()
	go ret.callExporter()

	return ret
}

func (p *Processor) Stop() {
	p.stopRequest <- true
}

func (p *Processor) Send(metric *Metric) {
	p.received <- metric
}

func (p *Processor) GetData() *MetricData {
	return p.lastExport
}

func (p *Processor) callExporter() {
	for {
		<-time.After(p.exportInterval)

		for _, h := range p.series {
			h.Calculate()
		}

		e := &MetricData{
			Metrics: p.metrics,
			Series:  p.series,
		}

		p.exportChan <- e
		p.lastExport = e

		p.metrics = make(map[string]*Metric)
		p.series = make(map[string]*Serie)
	}
}

func (p *Processor) process() {
	var m *Metric
	var e *MetricData

	for {
		select {
		case <-p.stopRequest:
			return

		case e = <-p.exportChan:
			p.export(e)

		case m = <-p.received:
			p.store(m)
		}
	}
}

func (p *Processor) store(m *Metric) {
	if m.Type == HistogramType || m.Type == SummaryType {
		if item, found := p.series[m.Name]; found {
			item.Increment(m.Value)
		} else {
			p.series[m.Name] = NewSerie(m.Name, m.Type, m.Tags, SecondResolution, m.Value)
		}
	} else {
		if item, found := p.metrics[m.Name]; found {
			item.Increment(m.Value)
		} else {
			p.metrics[m.Name] = NewMetric(m.Name, m.Type, m.Tags, m.Value)
		}
	}
}

func (p *Processor) export(e *MetricData) error {
	var err error
	if p.exporter != nil {
		err = p.exporter(e)

		if err != nil {
			log.Printf("fail to execute exporter: %s", err.Error)
		}
	}

	return err
}
