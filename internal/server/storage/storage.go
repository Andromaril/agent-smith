// Package storage содержит необходимое для работы с map как с хранилищем метрик
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/model"
)

// MemStorage хранит информацию о метриках для работы с хранилищем map
type MemStorage struct {
	Gauge     map[string]float64 // map c gauge метриками
	Counter   map[string]int64 // map c counter метриками
	WriteSync bool // для решения записи в файл метрик
	Path      string // путь, где лежит файл с метриками
	Mutex     *sync.Mutex
}

// Storage для работы с map и бд
type Storage interface {
	NewGauge(key string, value float64) error
	NewCounter(key string, value int64) error
	GetCounter(key string) (int64, error)
	GetGauge(key string) (float64, error)
	Load(file string) error
	Save(file string) error
	Init(path string, ctx context.Context) (*sql.DB, error)
	Ping() error
	GetIntMetric() (map[string]int64, error)
	GetFloatMetric() (map[string]float64, error)
	CounterAndGaugeUpdateMetrics(gauge []model.Gauge, counter []model.Counter) error
}

// Ping для работы с бд, пустой в данном пакет
func (m *MemStorage) Ping() error {
	return nil
}

// Init для работы с бд, пустой в данном пакет
func (m *MemStorage) Init(path string, ctx context.Context) (*sql.DB, error) {
	return nil, nil
}

// NewMemStorage для создания новых экземпляров MemStorage
func NewMemStorage(b bool, p string) *MemStorage {
	m := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), Path: p}
	m.SyncWrite(b)
	return &m
}

// SyncWrite от значений зависит, будет ли происходить запись в файл
func (m *MemStorage) SyncWrite(b bool) {
	m.WriteSync = b
}

// NewGauge для создание gauge-метрик
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

// NewCounter для создания counter метрик
func (m *MemStorage) NewCounter(key string, value int64) error {
	m.Counter[key] += value
	if m.WriteSync {
		err := m.Save(m.Path)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("not found %q", e.Error())
		}
	}
	return nil
}

// GetCounter для получения 1 counter метрики
func (m *MemStorage) GetCounter(key string) (int64, error) {
	k, ok := m.Counter[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
}

// GetGauge для получений 1 gauge-метрики
func (m *MemStorage) GetGauge(key string) (float64, error) {
	k, ok := m.Gauge[key]
	if !ok {
		return 0, fmt.Errorf("not found")
	}
	return k, nil
}

// Save для сохранения метрик в отдельный файл
func (m *MemStorage) Save(file string) error {
	// сериализуем структуру в JSON формат
	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("not found %q", e.Error())
	}
	return os.WriteFile(file, data, 0666)

}

// Load для чтения метрик из ранее созданного файла в Save
func (m *MemStorage) Load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("not found %q", e.Error())
	}
	json.Unmarshal(data, m)
	return nil
}

// GetIntMetric для получения 1 counter-метрики при работе с бд
func (m *MemStorage) GetIntMetric() (map[string]int64, error) {
	return m.Counter, nil
}

// GetFloatMetric для получения 1 gauge-метрики при работе с бд
func (m *MemStorage) GetFloatMetric() (map[string]float64, error) {
	return m.Gauge, nil
}

// CounterAndGaugeUpdateMetrics для получения списка gauge и counter метрик
func (m *MemStorage) CounterAndGaugeUpdateMetrics(gauge []model.Gauge, counter []model.Counter) error {
	for _, modelmetrics := range gauge {
		if err := m.NewGauge(modelmetrics.Key, modelmetrics.Value); err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("not found %q", e.Error())
		}
	}
	for _, modelmetrics := range counter {
		if err := m.NewCounter(modelmetrics.Key, modelmetrics.Value); err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("not found %q", e.Error())
		}
	}
	return nil
}
