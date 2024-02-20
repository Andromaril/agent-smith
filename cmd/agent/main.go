package main

import (
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"math/rand"
	"time"
)

// var RandomValue float64
// var PollCount int64


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
