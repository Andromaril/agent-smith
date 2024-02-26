package main

import (
	"math/rand"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
)

func UpdateMetric() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}

func main() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second*2)
		metric.SendAllMetricJSON()
		time.Sleep(time.Second*10)
		}
}
