package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/model"
)

type MemStorage struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	WriteSync bool
	Path      string
}

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

func (m *MemStorage) Ping() error {
	return nil
}
func (m *MemStorage) Init(path string, ctx context.Context) (*sql.DB, error) {
	return nil, nil
}

func NewMemStorage(b bool, p string) *MemStorage {
	m := MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), Path: p}
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
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("not found %q", e.Error())
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

func (m *MemStorage) Save(file string) error {
	// сериализуем структуру в JSON формат
	data, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("not found %q", e.Error())
	}
	return os.WriteFile(file, data, 0666)

}

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

func (m *MemStorage) GetIntMetric() (map[string]int64, error) {
	return m.Counter, nil
}

func (m *MemStorage) GetFloatMetric() (map[string]float64, error) {
	return m.Gauge, nil
}

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
