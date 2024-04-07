package creator

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var RandomValue float64
var PollCount int64

func CreateFloatMetric(metrics chan<- []model.Metrics) {
	for {
		modelmetrics := make([]model.Metrics, 0)
		var metr runtime.MemStats
		runtime.ReadMemStats(&metr)
		metric := map[string]float64{
			"Alloc":         float64(metr.Alloc),
			"BuckHashSys":   float64(metr.BuckHashSys),
			"Frees":         float64(metr.Frees),
			"GCCPUFraction": float64(metr.GCCPUFraction),
			"GCSys":         float64(metr.GCSys),
			"HeapAlloc":     float64(metr.HeapAlloc),
			"HeapIdle":      float64(metr.HeapIdle),
			"HeapInuse":     float64(metr.HeapInuse),
			"HeapObjects":   float64(metr.HeapObjects),
			"HeapReleased":  float64(metr.HeapReleased),
			"HeapSys":       float64(metr.HeapSys),
			"LastGC":        float64(metr.LastGC),
			"Lookups":       float64(metr.Lookups),
			"MCacheInuse":   float64(metr.MCacheInuse),
			"MCacheSys":     float64(metr.MCacheSys),
			"MSpanInuse":    float64(metr.MSpanInuse),
			"MSpanSys":      float64(metr.MSpanSys),
			"Mallocs":       float64(metr.Mallocs),
			"NextGC":        float64(metr.NextGC),
			"NumForcedGC":   float64(metr.NumForcedGC),
			"NumGC":         float64(metr.NumGC),
			"OtherSys":      float64(metr.OtherSys),
			"PauseTotalNs":  float64(metr.PauseTotalNs),
			"StackInuse":    float64(metr.StackInuse),
			"StackSys":      float64(metr.StackSys),
			"Sys":           float64(metr.Sys),
			"TotalAlloc":    float64(metr.TotalAlloc),
			"RandomValue":   float64(rand.Intn(1000)),
		}
		for name, metricvalue := range metric {
			value := metricvalue
			modelmetrics = append(modelmetrics, model.Metrics{ID: name, MType: "gauge", Value: &value})
		}
		delta := int64(1)
		modelmetrics = append(modelmetrics, model.Metrics{ID: "PollCount", MType: "counter", Delta: &delta})
		metrics <- modelmetrics
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}

}

func AddNewMetric(metrics chan<- []model.Metrics) {
	for {
		modelmetrics := make([]model.Metrics, 0)
		v, _ := mem.VirtualMemory()
		metric := map[string]float64{
			"TotalMemory": float64(v.Total),
			"FreeMemory":  float64(v.Free),
		}
		cpu, err := cpu.Percent(0, true)
		if err != nil {
			e := errormetric.NewMetricError(err)
			log.Printf("fatal get metric %q", e.Error())
		}
		for key, value := range cpu {
			metric[fmt.Sprintf("CPUutilization%d", key+1)] = float64(value)
		}
		for name, metricvalue := range metric {
			value := metricvalue
			modelmetrics = append(modelmetrics, model.Metrics{ID: name, MType: "gauge", Value: &value})
		}
		metrics <- modelmetrics
		time.Sleep(time.Second * time.Duration(flag.PollInterval))
	}
}
