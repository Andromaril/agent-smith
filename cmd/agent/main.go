package main

import (
	"math/rand"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
)

func main() {
	flag.ParseFlags()
	for {
		metric.PollCount++
		metric.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.ReportInterval))
		metric.SendAllMetric()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}
