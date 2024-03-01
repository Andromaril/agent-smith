package creator

import "runtime"

var RandomValue float64
var PollCount int64

func CreateFloatMetric() map[string]float64 {
	var metr runtime.MemStats
	runtime.ReadMemStats(&metr)
	return map[string]float64{
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
		"RandomValue":   RandomValue,
	}
}

func CreateIntMetric() map[string]int64 {
	return map[string]int64{
		"PollCount": PollCount,
	}
}
