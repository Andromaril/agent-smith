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
	//client := resty.New()
	// var t1 bool
	// var t2 bool
	// t1 = true
	// t2 = false
	//
	//
	for i = 0; ; i++ {
		time.Sleep(time.Second)
		//i++
		if i%flag.PollInterval == 0 {
			creator.PollCount++
			creator.RandomValue = rand.Float64()
			//i = i + flag.PollInterval
			// t2 = true
			// t1 = false
			// continue
		}
		if i%flag.ReportInterval == 0 {
			err := metric.SendAllMetricJSON2()
			if err != nil {
				panic(err)
			}
			// t1 = true
			// t2 = false
			// continue
			//i = i + flag.ReportInterval
		}
	}
}

