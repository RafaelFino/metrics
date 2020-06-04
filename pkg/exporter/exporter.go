package console

import (
	"encoding/json"
	"fmt"

	"github.com/rafaelfino/metrics/pkg/common"
)

type Exp struct {
}

func New() common.Exporter {
	return &Exp{}
}

func (e *Exp) Export(data *MetricData) error {
	raw, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(raw))
}
