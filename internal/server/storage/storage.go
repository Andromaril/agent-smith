package storage

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

func (m *MemStorage) NewGauge(key string, value float64) error{
	m.gauge[key] = value
	return nil
}

func (m *MemStorage) NewCounter(key string, value int64) error{
	m.counter[key] += value
	return nil
}
