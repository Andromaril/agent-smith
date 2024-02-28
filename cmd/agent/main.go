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
	var i int64
	client := resty.New()
	for {
		time.Sleep(time.Second)
		i++
		if i%flag.PollInterval == 0 {
			creator.PollCount++
			creator.RandomValue = rand.Float64()
			//i = i + flag.PollInterval
		}
		if i%flag.ReportInterval == 0 {
			err := metric.SendAllMetricJSON2(client)
			if err != nil {
				panic(err)
			}
			//i = i + flag.ReportInterval
		}
	}
	//time.Sleep(100*time.Second)
	//client := resty.New()
	//go func() {
	//for {
	//err := metric.SendAllMetricJSON()
	//if err != nil {
	//panic(err)
	//}
	//time.Sleep(time.Second * time.Duration(flag.ReportInterval))
	//}
	//}()
	//for {
	//UpdateMetric()
	//}
}
