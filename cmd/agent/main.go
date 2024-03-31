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

func UpdateMetric() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}

func main() {
	flag.ParseFlags()
	//var i int64
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
	wg.Add(1)
	// for i = 0; ; i++ {
	// 	time.Sleep(time.Second)
	// 	if i%flag.PollInterval == 0 {
	// 		creator.PollCount++
	// 		creator.RandomValue = rand.Float64()
	// 		creator.CreateFloatMetric(storage)
	// 		creator.CreateIntMetric(storage)

	// 	}
	// 	if i%flag.ReportInterval == 0 {
	// 		operation := func() error {
	// 			err := metric.SendAllMetricJSON(sugar, storage)
	// 			return err
	// 		}
	// 		err2 := retry.Retry(operation)

	// 		if err2 != nil {
	// 			sugar.Errorw(
	// 				"error when send mentric")
	// 		}

	// 	}
	// }
	go func() {
		defer wg.Done()
		for {
			func() {
				creator.PollCount++
				creator.RandomValue = rand.Float64()
				creator.CreateFloatMetric(storage)
				creator.CreateIntMetric(storage)
			}()
			//time.Sleep(time.Duration(flag.PollInterval) * time.Second)
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
			//time.Sleep(time.Duration(flag.PollInterval) * time.Second)
			time.Sleep(time.Duration(flag.ReportInterval) * time.Second)
		}
	}()
	wg.Wait()
}
