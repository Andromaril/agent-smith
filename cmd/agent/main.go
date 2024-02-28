package main

import (
	"math/rand"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/go-resty/resty/v2"
)

func UpdateMetric() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}

func main() {
	flag.ParseFlags()
	time.Sleep(100*time.Second)
	client := resty.New()
	go func() {
		for {
			err := metric.SendAllMetricJSON(client)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second * time.Duration(flag.ReportInterval))
		}
	}()
	for {
		UpdateMetric()
	}
}
