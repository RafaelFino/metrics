package histogram

import (
	"time"

	"github.com/rafaelfino/metrics/pkg/common"
)

type Histogram struct {
	name       string
	createdAt  time.Time
	lastAt     time.Time
	tags       map[string]string
	data       map[int64]float64
	resolution common.HistogramResolution
}

func New(name string, tags map[string]string, resolution common.HistogramResolution, value float64) common.Series {
	return &Histogram{
		name:       name,
		createdAt:  time.Now(),
		lastAt:     time.Now(),
		tags:       tags,
		data:       map[int64]float64{getKey(resolution): value},
		resolution: resolution,
	}
}

func getKey(resolution common.HistogramResolution) int64 {
	when := time.Now()

	switch resolution {
	case common.SecondResolution:
		when = when.Truncate(time.Second)
	case common.MinuteResolution:
		when = when.Truncate(time.Minute)
	case common.HourResolution:
		when = when.Truncate(time.Hour)
	case common.DayResolution:
		when = when.Truncate(time.Hour * 24)
	}

	return when.Unix()
}

func (h *Histogram) Clear() {
	h.data = make(map[int64]float64)
	h.createdAt = time.Now()
	h.lastAt = time.Now()
}

func (h *Histogram) CreatedAt() time.Time {
	return h.createdAt
}

func (h *Histogram) LastAt() time.Time {
	return h.lastAt
}

func (h *Histogram) Tags() map[string]string {
	return h.tags
}

func (h *Histogram) Name() string {
	return h.name
}

func (h *Histogram) Increment(value float64) {
	h.data[getKey(h.resolution)] += value
	h.lastAt = time.Now()
}

func (h *Histogram) Count() int {
	return len(h.data)
}

func (h *Histogram) Sum() float64 {
	ret := float64(0)

	for _, i := range h.data {
		ret += i
	}

	return ret
}

func (h *Histogram) Avg() float64 {
	c := h.Count()

	if c == 0 {
		return 0
	}

	return h.Sum() / float64(c)
}

func (h *Histogram) Data() []float64 {
	ret := make([]float64, len(h.data))

	p := 0
	for _, m := range h.data {
		ret[p] = m
		p++
	}

	return ret
}
