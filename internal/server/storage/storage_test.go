package storage

import (
	"math/rand"
	"testing"

	"github.com/andromaril/agent-smith/internal/model"
	"github.com/sirupsen/logrus"
)

func TestMemStorage_NewGauge(t *testing.T) {
	type args struct {
		gauge map[string]float64
	}
	tests := []struct {
		args args
		name string
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
			s := NewMemStorage(false, "test")
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
		args args
		name string
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
			s := NewMemStorage(false, "test")
			for key, value := range tt.args.counter {
				err := s.NewCounter(key, value)
				if err != nil {
					t.Errorf("NewCounter error = %v", err)
				}
			}
		})
	}
}

func TestMemStorage_CounterAndGaugeUpdateMetrics(t *testing.T) {
	type args struct {
		gauge   []model.Gauge
		counter []model.Counter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				gauge: []model.Gauge{
					{
						Key:   "Gauge",
						Value: 2.2,
					},
				},
				counter: []model.Counter{
					{
						Key:   "Counter",
						Value: 2,
					},
				},
			},
		},
		{
			name: "test2",
			args: args{
				gauge: []model.Gauge{
					{
						Key:   "Gauge",
						Value: 2.4,
					},
				},
				counter: []model.Counter{
					{
						Key:   "Counter",
						Value: 4,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage(false, "test")
			err := s.CounterAndGaugeUpdateMetrics(tt.args.gauge, tt.args.counter)
			if err != nil {
				t.Errorf("NewGauge error = %v", err)
			}
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test",
			args: args{key: "test"},
			want: 2.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage(false, "test")
			m := s.Gauge
			m["test"] = 2.2
			got, _ := s.GetGauge(tt.args.key)
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "test",
			args: args{key: "test"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage(false, "test")
			m := s.Counter
			m["test"] = 2
			got, _ := s.GetCounter(tt.args.key)
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Save(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			path: "/tmp/test-1.json",
			args: args{
				file: "metric gauge, counter",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage(false, tt.path)
			if err := s.Save(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_Load(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			path: "/tmp/test-1.json",
			args: args{
				file: "metric gauge, counter",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemStorage(false, tt.path)
			if err := s.Load(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkAddGauge(b *testing.B) {
	s := NewMemStorage(false, "test")
	for i := 0; i < b.N; i++ {
		key := "Gauge"
		val := rand.Float64()
		err := s.NewGauge(key, val)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func BenchmarkAddCounter(b *testing.B) {
	s := NewMemStorage(false, "test")
	for i := 0; i < b.N; i++ {
		key := "counter"
		val := rand.Int63()
		err := s.NewCounter(key, val)
		if err != nil {
			logrus.Error(err)
		}
	}
}
