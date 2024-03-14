package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemStorage struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	WriteSync bool
	Path      string
}

func NewMemStorage(b bool, p string) *MemStorage {
	m := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), Path: p}
	// return &MemStorage{
	// 	Gauge:   make(map[string]float64),
	// 	Counter: make(map[string]int64),
	// 	writeSync:
	// }
	m.SyncWrite(b)
	return &m
}

func (m *MemStorage) SyncWrite(b bool) {
	m.WriteSync = b
}

func (m *MemStorage) NewGauge(key string, value float64) error {
	m.Gauge[key] = value
	if m.WriteSync {
		err := m.Save(m.Path)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (m *MemStorage) NewCounter(key string, value int64) error {
	m.Counter[key] += value
	if m.WriteSync {
		err := m.Save(m.Path)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (m *MemStorage) GetCounter(key string) (int64, error) {
	k, ok := m.Counter[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
}

func (m *MemStorage) GetGauge(key string) (float64, error) {
	k, ok := m.Gauge[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
}

func (m *MemStorage) PrintMetric() string {
	var result string
	for k1, v1 := range m.Gauge {
		result += fmt.Sprintf("%s: %v\n", k1, v1)
	}
	for k2, v2 := range m.Counter {
		result += fmt.Sprintf("%s: %v\n", k2, v2)
	}
	return result
}

func (m *MemStorage) Save(file string) error {
	// сериализуем структуру в JSON формат
	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0666)

}

func (m *MemStorage) Load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	json.Unmarshal(data, m)
	return nil
}

func (m *MemStorage) GetIntMetric() map[string]int64 {
	return m.Counter
}

func (m *MemStorage) GetFloatMetric() map[string]float64 {
	return m.Gauge
}
