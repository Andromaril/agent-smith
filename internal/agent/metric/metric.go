package metric

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/andromaril/agent-smith/internal/flag"
	"runtime"
)

var RandomValue float64
var PollCount int64


func SendGaugeMetric(name string, value float64) {
	client := resty.New()
	url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%v", name, value)
	fmt.Print(url)
	_, err := client.R().Post(url)
	if err != nil {
		panic(err)
	}
}

func SendCounterMetric(name string, value int64) {
	client := resty.New()
	url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%v", name, value)
	_, err := client.R().Post(url)
	if err != nil {
		panic(err)
	}
}

func SendAllMetric() {
	var metr runtime.MemStats
	runtime.ReadMemStats(&metr)
	SendGaugeMetric("Alloc", float64(metr.Alloc))
	SendGaugeMetric("BuckHashSys", float64(metr.BuckHashSys))
	SendGaugeMetric("Frees", float64(metr.Frees))
	SendGaugeMetric("GCCPUFraction", float64(metr.GCCPUFraction))
	SendGaugeMetric("GCSys", float64(metr.GCSys))
	SendGaugeMetric("HeapAlloc", float64(metr.HeapAlloc))
	SendGaugeMetric("HeapIdle", float64(metr.HeapIdle))
	SendGaugeMetric("HeapInuse", float64(metr.HeapInuse))
	SendGaugeMetric("HeapObjects", float64(metr.HeapObjects))
	SendGaugeMetric("HeapReleased", float64(metr.HeapReleased))
	SendGaugeMetric("HeapSys", float64(metr.HeapSys))
	SendGaugeMetric("LastGC", float64(metr.LastGC))
	SendGaugeMetric("Lookups", float64(metr.Lookups))
	SendGaugeMetric("CacheInuse", float64(metr.MCacheInuse))
	SendGaugeMetric("MCacheSys", float64(metr.MCacheSys))
	SendGaugeMetric("MSpanInuse", float64(metr.MSpanInuse))
	SendGaugeMetric("MSpanSys", float64(metr.MSpanSys))
	SendGaugeMetric("Mallocs", float64(metr.Mallocs))
	SendGaugeMetric("NextGC", float64(metr.NextGC))
	SendGaugeMetric("NumForcedGC", float64(metr.NumForcedGC))
	SendGaugeMetric("NumGC", float64(metr.NumGC))
	SendGaugeMetric("OtherSys", float64(metr.OtherSys))
	SendGaugeMetric("PauseTotalNs", float64(metr.PauseTotalNs))
	SendGaugeMetric("StackInuse", float64(metr.StackInuse))
	SendGaugeMetric("StackSys", float64(metr.StackSys))
	SendGaugeMetric("Sys", float64(metr.Sys))
	SendGaugeMetric("TotalAlloc", float64(metr.TotalAlloc))
	SendGaugeMetric("RandomValue", RandomValue)
	SendCounterMetric("PollCount", PollCount)

}
