package handlergrpc

import (
	"context"
	"database/sql"

	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	pb "github.com/andromaril/agent-smith/pkg/proto"
)

// MetricServer поддерживает все необходимые методы сервера.
type MetricServer struct {
	pb.UnimplementedMetricServer
	db       storagedb.Interface
	database *sql.DB
}

func NewMetricServer(i storagedb.Interface, db *sql.DB) *MetricServer {
	return &MetricServer{db: i, database: db}
}

func (s *MetricServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse
	gauge := make([]model.Gauge, 0)
	counter := make([]model.Counter, 0)
	for _, v := range in.Counter {
		counter = append(counter, model.Counter{Key: v.Key, Value: v.Value})
	}
	for _, v := range in.Gauge {
		gauge = append(gauge, model.Gauge{Key: v.Key, Value: v.Value})
	}

	if err := s.db.CounterAndGaugeUpdateMetrics(gauge, counter); err != nil {
		return &response, err
	}
	return &response, nil
}
