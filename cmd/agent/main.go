package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/retry"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

// func UpdateMetric() {
// 	for {
// 		creator.PollCount++
// 		creator.RandomValue = rand.Float64()
// 		time.Sleep(time.Second * time.Duration(flag.PollInterval))
// 	}
// }

func main() {
	flag.ParseFlags()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow(
		"Starting agent")
	storage := storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}}

	var wg sync.WaitGroup
	ratelimit := flag.RateLimit
	wg.Add(int(ratelimit))

	go func() {
		defer wg.Done()
		for {
			func() {
				creator.PollCount++
				creator.RandomValue = rand.Float64()
				creator.CreateFloatMetric(storage)
				creator.CreateIntMetric(storage)
				storage.AddNewMetric()
			}()
			time.Sleep(time.Duration(flag.PollInterval) * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			func() {
				operation := func() error {
					err := metric.SendAllMetricJSON(sugar, storage)
					return err
				}
				err2 := retry.Retry(operation)

				if err2 != nil {
					sugar.Errorw(
						"error when send mentric")
				}
			}()
			time.Sleep(time.Duration(flag.ReportInterval) * time.Second)
		}
	}()
	wg.Wait()
}
