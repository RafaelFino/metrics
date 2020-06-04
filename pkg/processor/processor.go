package processor

import (
	"log"
	"time"

	. "github.com/rafaelfino/metrics/pkg/common"
)

type Processor struct {
	received    chan *Metric
	exportChan  chan *MetricData
	stopRequest chan bool

	counters   map[string]*Metric
	gauges     map[string]*Metric
	histograms map[string]*Serie

	exporter       Exporter
	exportInterval time.Duration
}

func New(exportInterval time.Duration, exporter Exporter) *Processor {
	ret := &Processor{
		received:       make(chan *Metric, 256),
		exportChan:     make(chan *MetricData, 10),
		stopRequest:    make(chan bool),
		exporter:       exporter,
		exportInterval: exportInterval,
		counters:       make(map[string]*Metric),
		gauges:         make(map[string]*Metric),
		histograms:     make(map[string]*Serie),
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

func (p *Processor) callExporter() {
	for {
		<-time.After(p.exportInterval)
		p.exportChan <- &MetricData{
			Counters:   p.counters,
			Gauges:     p.gauges,
			Histograms: p.histograms,
		}
	}
}

func (p *Processor) process() {
	var m *Metric
	var e *MetricData

	for {
		select {
		case <-p.stopRequest:
			log.Printf("calling export and exit...")
			return

		case e = <-p.exportChan:
			log.Printf("calling export...")
			p.export(e)
			p.clear()

		case m = <-p.received:
			log.Printf("calling store...")
			p.store(m)
		}
	}
}

func (p *Processor) store(m *Metric) {
	switch m.Type {
	case CounterType:
		if item, found := p.counters[m.Name]; found {
			item.Value += m.Value
		} else {
			p.counters[m.Name] = m
		}
	case GaugeType:
		p.gauges[m.Name] = m
	case HistogramType:
		if item, found := p.histograms[m.Name]; found {
			item.Increment(m.Value)
		} else {
			p.histograms[m.Name] = NewSerie(m.Name, HistogramType, m.Tags, SecondResolution, m.Value)
		}
	}
}

func (p *Processor) clear() {
	p.counters = make(map[string]*Metric)
	p.gauges = make(map[string]*Metric)
	p.histograms = make(map[string]*Serie)
}

func (p *Processor) export(e *MetricData) error {
	log.Printf("start export...")
	var err error
	if p.exporter != nil {
		err = p.exporter.Export(e)

		if err != nil {
			log.Printf("fail to execute exporter: %s", err.Error)
		}
	}

	return err
}
