package storage

import "fmt"

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *MemStorage) NewGauge(key string, value float64) error {
	m.gauge[key] = value
	return nil
}

func (m *MemStorage) NewCounter(key string, value int64) error {
	m.counter[key] += value
	return nil
}

func (m *MemStorage) GetCounter(key string) int64 {
	k := m.counter[key]
	return k
}

func (m *MemStorage) GetGauge(key string) float64 {
	k := m.gauge[key]
	return k
}

func (m *MemStorage) PrintMetric() string {
	var result string
	for k1, v1 := range m.gauge {
		result += fmt.Sprintf("%s: %v\n", k1, v1)
	}
	for k2, v2 := range m.counter {
		result += fmt.Sprintf("%s: %v\n", k2, v2)
	}
	return result
}