package main

import (
	"math/rand"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func UpdateMetric() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}

func main() {
	flag.ParseFlags()
	var i int64
	logger, err1 := zap.NewDevelopment()
	if err1 != nil {
		panic(err1)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow(
		"Starting agent")
	storage := storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}
	for i = 0; ; i++ {
		time.Sleep(time.Second)
		if i%flag.PollInterval == 0 {
			creator.PollCount++
			creator.RandomValue = rand.Float64()
			creator.CreateFloatMetric(storage)
			creator.CreateIntMetric(storage)

		}
		if i%flag.ReportInterval == 0 {
			err := metric.SendAllMetricJSON(storage)
			if err != nil {
				sugar.Errorw(
					"Error send metric", err)
			}
		}
	}
}
