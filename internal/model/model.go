// Package model хранит модели используемые агентом и сервисом
package model

// Metrics хранит информацию о метриках
type Metrics struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

// Counter хранит информацию о метриках формата int64
type Counter struct {
	Key   string
	Value int64
}

// Gauge хранит информацию о метриках формата float64
type Gauge struct {
	Key   string
	Value float64
}
