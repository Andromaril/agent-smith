package main

import (
	"math/rand"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
)

//func NewMetric() {
//for {
//metric.SendAllMetric()
//time.Sleep(time.Second * time.Duration(flag.ReportInterval))
//}
//}

func UpdateMetric() {
	for {
		creator.PollCount++
		creator.RandomValue = rand.Float64()
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}

func main() {
	flag.ParseFlags()
	go UpdateMetric()
	time.Sleep(time.Second)
	for {
		err := metric.SendAllMetric()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * time.Duration(flag.ReportInterval))
	}
}
