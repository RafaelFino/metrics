package console

import (
	"encoding/json"
	"fmt"

	"github.com/rafaelfino/metrics"
)

type Exp struct {
}

func New(config map[string]string) metrics.Exporter {
	return &Exp{}
}

func (e *Exp) Export(data *metrics.MetricData) error {
	raw, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(raw))
}
