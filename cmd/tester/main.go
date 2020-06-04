package main

import (
	"log"
	"time"

	"github.com/rafaelfino/metrics"
	"github.com/rafaelfino/metrics/exporters/console"
)

func main() {
	log.Println("Creating processor")
	exp := console.New(map[string]string{})
	p := metrics.Processor.New(time.Second*10, exp)
	defer p.Stop()

	for {
		log.Println("Sending metrics...")

		p.Send(&metrics.Metric{Name: "counter.click", Tags: map[string]string{"tag1": "value1", "tag2": "value2"}, When: time.Now(), Type: metrics.Counter})

		time.Sleep(time.Second)
	}
}
