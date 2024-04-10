package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/retry"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func worker(jobs <-chan []model.Metrics, sugar zap.SugaredLogger) {
	for j := range jobs {
		operation := func() error {
			err := metric.SendMetricJSON(sugar, j)
			return err
		}
		err2 := retry.Retry(operation)

		if err2 != nil {
			sugar.Errorw(
				"error when send mentric")
		}
		time.Sleep(time.Second * time.Duration(flag.ReportInterval))
	}
}

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
	var wg sync.WaitGroup
	ratelimit := flag.RateLimit
	jobs := make(chan []model.Metrics, runtime.GOMAXPROCS(0))
	wg.Add(ratelimit)
	go creator.CreateFloatMetric(jobs)
	go creator.AddNewMetric(jobs)
	defer close(jobs)
	for w := 1; w <= ratelimit; w++ {
		go worker(jobs, sugar)
	}
	wg.Wait()
}
