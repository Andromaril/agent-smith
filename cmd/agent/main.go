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
	time.Sleep(time.Second)
	//go UpdateMetric()
	//for {
	//err := metric.SendAllMetricJSON()
	//if err != nil {
	//panic(err)
	//}
	//time.Sleep(time.Second * time.Duration(flag.ReportInterval))
	//}
	go func() {
		for {
			err := metric.SendAllMetricJSON()
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
