package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

type SerializedMemStorage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m *MemStorage) SetMetricsData(gauge map[string]float64, counter map[string]int64) {

	m.Gauge = gauge
	m.Counter = counter
}

func (m *MemStorage) NewGauge(key string, value float64) error {
	m.Gauge[key] = value
	return nil
}

func (m *MemStorage) NewCounter(key string, value int64) error {
	m.Counter[key] += value
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
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// сохраняем данные в файл
	if err := os.WriteFile(file, data, 0666); err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) Load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// data2 := &SerializedMemStorage{}
	// if err := json.Unmarshal(data, &data2); err != nil {
	// 	return err
	// }
	// m.SetMetricsData(data2.Gauge, data2.Counter)
	json.Unmarshal(data, m)
	return nil
}

// func RestoreData(m *MemStorage, value bool) {
// 	if value {
// 		m.Load()
// 	}
// }
