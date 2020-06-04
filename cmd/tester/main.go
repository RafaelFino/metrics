package main

import (
	"log"
	"time"

	"github.com/rafaelfino/metrics/pkg/common"
	"github.com/rafaelfino/pkg/exporter/console"
	"github.com/rafaelfino/pkg/processor"
)

func main() {
	log.Println("Creating processor")

	exp := console.New()

	p := processor.New(time.Second*10, exp)
	defer p.Stop()

	for {
		log.Println("Sending metrics...")

		p.Send(&common.Metric{Name: "counter.click", Tags: map[string]string{"tag1": "value1", "tag2": "value2"}, When: time.Now(), Type: common.CounterType})

		time.Sleep(time.Second)
	}
}
