package metrics

import "time"

type HistogramResolution string

const (
	Second HistogramResolution = "second"
	Minute HistogramResolution = "minute"
	Hour   HistogramResolution = "hour"
	Day    HistogramResolution = "day"
)

type Histogram struct {
	name       string
	createdAt  time.Time
	lastAt     time.Time
	tags       map[string]string
	data       map[int64]float64
	resolution HistogramResolution
}

func New(name string, tags map[string]string, resolution HistogramResolution, value float64) *Histogram {
	return &Histogram{
		name:       name,
		createdAt:  time.Now(),
		lastAt:     time.Now(),
		tags:       tags,
		data:       map[int64]float64{getKey(resolution): value},
		resolution: resolution,
	}
}

func getKey(resolution HistogramResolution) int64 {
	when := time.Now()

	switch resolution {
	case Second:
		when = when.Truncate(time.Second)
	case Minute:
		when = when.Truncate(time.Minute)
	case Hour:
		when = when.Truncate(time.Hour)
	case Day:
		when = when.Truncate(time.Hour * 24)
	}

	return when.Unix()
}

func (h *Histogram) Clear() {
	h.data = make(map[int64]floate64)
	h.createdAt = time.Now
	h.lastAt = time.Now
}

func (h *Histogram) CreatedAt() time.Time {
	return s.createdAt
}

func (h *Histogram) Last() time.Time {
	return s.lastAt
}

func (h *Histogram) Tags() map[string]string {
	return s.tags
}

func (h *Histogram) Name() string {
	return s.name
}

func (h *Histogram) Increment(value float64) {
	s.data[getKey(h.resolution)] += value
	s.last = last
}

func (h *Histogram) Count() int {
	return len(s.data)
}

func (h *Histogram) Sum() float64 {
	ret := float64(0)

	for _, i := range s.data {
		ret += i
	}

	return ret
}

func (h *Histogram) Avg() float64 {
	c := s.Count()

	if c == 0 {
		return 0
	}

	return s.Sum / c
}

func (h *Histogram) Data() []float64 {
	ret := make([]float64, len(s.data))

	p := 0
	for _, m := range s.data {
		ret[p] = m
		p++
	}

	return ret
}
