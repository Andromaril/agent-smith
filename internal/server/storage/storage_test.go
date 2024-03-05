package storage

import (
	"testing"
)

func TestMemStorage_NewGauge(t *testing.T) {
	type args struct {
		gauge map[string]float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				gauge: map[string]float64{
					"gauge1": 1.1,
					"gauge2": 2.2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage()
			for key, value := range tt.args.gauge {
				err := s.NewGauge(key, value)
				if err != nil {
					t.Errorf("NewGauge error = %v", err)
				}
			}
		})
	}
}

func TestMemStorage_NewCounter(t *testing.T) {
	type args struct {
		counter map[string]int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				counter: map[string]int64{
					"counter1": 3,
					"counter2": 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage()
			for key, value := range tt.args.counter {
				err := s.NewCounter(key, value)
				if err != nil {
					t.Errorf("NewCounter error = %v", err)
				}
			}
		})
	}
}
