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
	flag.ParseFlags()
	var i int64
	var t1 bool
	var t2 bool
	t1 = true
	t2 = false
	for i = 0; ; i++ {
		time.Sleep(time.Second)
		if t1 && i%flag.PollInterval == 0 {
			creator.PollCount++
			creator.RandomValue = rand.Float64()
			t2 = true
			t1 = false
			time.Sleep(time.Second * time.Duration(flag.PollInterval))

		}
		if t2 && i%flag.ReportInterval == 0 {
			err := metric.SendAllMetricJSON2()
			if err != nil {
				panic(err)
			}
			t1 = true
			t2 = false
			time.Sleep(time.Second * time.Duration(flag.ReportInterval))
		}
	}
}
