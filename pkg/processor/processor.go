package processor

import (
	"log"
	"sync"
	"time"

	. "github.com/rafaelfino/metrics/metrics/common"
	. "github.com/rafaelfino/metrics/metrics/histogram"
	. "github.com/rafaelfino/metrics/metrics/summary"
	"github.com/rafaelfino/metrics/pkg/histogram"
	"github.com/rafaelfino/metrics/pkg/summary"
)

type Processor struct {
	received    chan *Metric
	stopRequest chan bool

	counters   map[string]*Metric
	gauges     map[string]*Metric
	histograms map[string]*Series
	summaries  map[string]*Series

	exporter       Exporter
	ExportInterval time.Duration
	mu             *sync.Mutex
}

func New(exportInterval time.Duration, exporter Exporter) *Processor {
	ret := &Processor{
		running:        true,
		received:       make(chan *Metric, 256),
		stopRequest:    make(chan bool),
		exporter:       exporter,
		exportInterval: exportInterval,
		counters:       make(map[string]*Metric),
		gauges:         make(map[string]*Metric),
		histograms:     make(map[string]*Series),
		summaries:      make(map[string]*Series),
		mutex:          &sync.Mutex{},
	}

	go ret.process()

	return ret
}

func (p *Processor) Stop() {
	p.stopRequest <- true
}

func (p *Processor) Send(metric *Metric) {
	p.received <- metric
}

func (p *Processor) process() {
	for {
		select {
		case m <- p.received:
			p.store(m)
		case <-p.stopRequest:
			p.Export()
			return
		case <-time.After(p.exportInterval):
			p.Export()
		}
	}
}

func (p *Processor) store(m *Metric) {
	switch m.Type {
	case CounterType:
		if item, found := p.counters[m.Name]; found {
			p.counters[m.Name].Value += m.Value
		} else {
			p.counters[m.Name] = m
		}
	case GaugeType:
		p.gauges[m.Name] = m
	case HistogramType:
		if item, found := p.histograms[m.Name]; !found {
			p.histograms[m.Name] = histogram.New(m.Name, m.Tags, Second, m.Value)
		} else {
			item.Increment(m.Value)
		}
	case SummaryType:
		if item, found := p.summaries[m.Name]; !found {
			p.summaries[m.Name] = summary.New(m.Name, m.Tags, m.Value)
		} else {
			item.Increment(m.Value)
		}
	}
}

func (p *Processor) clear() {
	p.counters = make(map[string]*Metric)
	p.gauges = make(map[string]*Metric)
	p.histograms = make(map[string]*Series)
	p.summaries = make(map[string]*Series)
}

func (p *Processor) Export() error {
	var err error
	if p.export != nil {
		err = p.export.Export(&MetricData{
			Counters:   p.counters,
			Gauges:     p.gauges,
			Histograms: p.histograms,
			Summaries:  p.summaries,
		})

		if err != nil {
			log.Printf("fail to execute exporter: %s", err.Error)
		}
	}
	p.clear()

	return err
}
