package main

import (
	"time"

	"github.com/rafaelfino/metrics"
)

func main() {
	exp := metrics.NewConsoleExporter()

	p := metrics.NewMetricProcessor(time.Second*10, exp)
	defer p.Stop()

	wait := time.Second / 1000

	tags := map[string]string{"tag1": "value1", "tag2": "value2"}

	for i := 0; i < 50000; i++ {
		p.Send(metrics.NewMetric("counter.fixed", metrics.CounterType, tags, 1))
		p.Send(metrics.NewMetric("counter.var", metrics.CounterType, tags, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("gauger.fixed", metrics.GaugeType, tags, 5))
		p.Send(metrics.NewMetric("gauger.var", metrics.GaugeType, tags, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("histogram.fixed", metrics.HistogramType, tags, 2))
		p.Send(metrics.NewMetric("histogram.var", metrics.HistogramType, tags, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("summary.fixed", metrics.SummaryType, tags, 2))
		p.Send(metrics.NewMetric("summary.var", metrics.SummaryType, tags, float64(time.Now().Unix()%10)))

		time.Sleep(wait)
	}
}
