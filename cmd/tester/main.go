package main

import (
	"log"
	"time"

	"github.com/rafaelfino/metrics/pkg/common"
	"github.com/rafaelfino/metrics/pkg/exporter/console"
	"github.com/rafaelfino/metrics/pkg/processor"
)

func main() {
	log.Println("Creating processor")

	exp := console.New()

	p := processor.New(time.Second*5, exp)
	defer p.Stop()

	wait := time.Second

	for {
		log.Println("Sending metrics...")

		p.Send(&common.Metric{Name: "counter.click", Tags: map[string]string{"tag1": "value1", "tag2": "value2"}, When: time.Now(), Type: common.CounterType, Value: 1})
		p.Send(&common.Metric{Name: "gauge.click", Tags: map[string]string{"tag1": "value1", "tag2": "value2"}, When: time.Now(), Type: common.GaugeType, Value: float64(time.Now().Unix())})
		p.Send(&common.Metric{Name: "histogram.click", Tags: map[string]string{"tag1": "value1", "tag2": "value2"}, When: time.Now(), Type: common.HistogramType, Value: float64(time.Now().Unix() % 3)})

		time.Sleep(wait)

		//wait *= 2
	}

	log.Println("Stop processor")
}
