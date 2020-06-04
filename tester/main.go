package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rafaelfino/metrics"
)

func main() {
	p := metrics.NewMetricProcessor(time.Second*2, ConsoleExport)
	defer p.Stop()

	wait := time.Second / 1000
	start := time.Now()

	tags := map[string]string{"tag1": "value1", "tag2": "value2"}

	for time.Since(start) < (time.Second * 10) {
		p.Send(metrics.NewMetric("counter.fixed", metrics.CounterType, tags, 1))
		p.Send(metrics.NewMetric("counter.var", metrics.CounterType, tags, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("gauger.fixed", metrics.GaugeType, nil, 5))
		p.Send(metrics.NewMetric("gauger.var", metrics.GaugeType, nil, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("histogram.fixed", metrics.HistogramType, tags, 2))
		p.Send(metrics.NewMetric("histogram.var", metrics.HistogramType, tags, float64(time.Now().Unix()%10)))

		p.Send(metrics.NewMetric("summary.fixed", metrics.SummaryType, tags, 2))
		p.Send(metrics.NewMetric("summary.var", metrics.SummaryType, tags, float64(time.Now().Unix()%10)))

		time.Sleep(wait)
	}
}

func ConsoleExport(data *metrics.MetricData) error {
	raw, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(raw))

	return err
}
