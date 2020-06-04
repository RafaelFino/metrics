package summary

import (
	"time"

	"github.com/rafaelfino/metrics/pkg/common"
)

type Summary struct {
	name      string
	createdAt time.Time
	lastAt    time.Time
	tags      map[string]string
	data      []float64
}

func New(name string, tags map[string]string, value float64) common.Series {
	return &Summary{
		name:      name,
		createdAt: time.Now(),
		lastAt:    time.Now(),
		tags:      tags,
		data:      []float64{value},
	}
}

func (s *Summary) Clear() {
	s.data = []float64{}
	s.createdAt = time.Now()
	s.lastAt = time.Now()
}

func (s *Summary) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Summary) LastAt() time.Time {
	return s.lastAt
}

func (s *Summary) Tags() map[string]string {
	return s.tags
}

func (s *Summary) Name() string {
	return s.name
}

func (s *Summary) Increment(value float64) {
	last := time.Now()
	s.data = append(s.data, value)
	s.lastAt = last
}

func (s *Summary) Count() int {
	return len(s.data)
}

func (s *Summary) Sum() float64 {
	ret := float64(0)

	for _, i := range s.data {
		ret += i
	}

	return ret
}

func (s *Summary) Avg() float64 {
	c := s.Count()

	if c == 0 {
		return 0
	}

	return s.Sum() / float64(c)
}

func (s *Summary) Data() []float64 {
	return s.data
}
