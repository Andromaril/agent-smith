package main

import (
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"math/rand"
	"time"
)

// var RandomValue float64
// var PollCount int64

func main() {
	for {
		metric.PollCount++
		metric.RandomValue = rand.Float64()
		time.Sleep(time.Second*2)
		metric.SendAllMetric()
		time.Sleep(time.Second*10)
		}
}

