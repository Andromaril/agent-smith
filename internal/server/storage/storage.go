package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andromaril/agent-smith/internal/flag"
)

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

func (m *MemStorage) GetCounter(key string) (int64, error) {
	k, ok := m.counter[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
}

func (m *MemStorage) GetGauge(key string) (float64, error) {
	k, ok := m.gauge[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
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

func Save(m *MemStorage) error {
	// сериализуем структуру в JSON формат
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// сохраняем данные в файл
	if err := os.WriteFile(flag.FileStoragePath, data, 0666); err != nil {
		return err
	}
	return nil
}

func Load(m *MemStorage) error {
	data, err := os.ReadFile(flag.FileStoragePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	return nil
}

func RestoreData(m *MemStorage) {
	if flag.Restore {
		Load(m)
	}
}
