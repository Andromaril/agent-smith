// Package model хранит модели используемые агентом и сервисом
package model

// Metrics для описания метрик
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Counter для описания метрик формата int64
type Counter struct {
	Key   string
	Value int64
}

// Gauge для описания метрик формата float64
type Gauge struct {
	Key   string
	Value float64
}
