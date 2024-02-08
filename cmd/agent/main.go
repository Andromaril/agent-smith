package main

import (
	"flag"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"math/rand"
	"time"
)

// var RandomValue float64
// var PollCount int64


var reportInterval = flag.Int64("r", 10, "time to sleep for report interval")
var pollInterval = flag.Int64("p", 2, "time to sleep for poll interval")

func main() {
	flag.Parse();
	for {
		metric.PollCount++
		metric.RandomValue = rand.Float64()
		time1 := rand.Int63n(*pollInterval)
		time.Sleep(time.Second * time.Duration(time1))
		metric.SendAllMetric()
		time2 := rand.Int63n(*reportInterval)
		time.Sleep(time.Second * time.Duration(time2))
	}
}
