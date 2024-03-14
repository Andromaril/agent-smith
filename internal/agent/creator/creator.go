package creator

import (
	"runtime"

	"github.com/andromaril/agent-smith/internal/server/storage"
)

var RandomValue float64
var PollCount int64

func CreateFloatMetric(storage storage.MemStorage) {
	var metr runtime.MemStats
	runtime.ReadMemStats(&metr)
	//return map[string]float64{
	storage.NewGauge("Alloc", float64(metr.Alloc))
	storage.NewGauge("BuckHashSys", float64(metr.BuckHashSys))
	storage.NewGauge("Frees", float64(metr.Frees))
	storage.NewGauge("GCCPUFraction", float64(metr.GCCPUFraction))
	storage.NewGauge("GCSys", float64(metr.GCSys))
	storage.NewGauge("HeapAlloc", float64(metr.HeapAlloc))
	storage.NewGauge("HeapIdle", float64(metr.HeapIdle))
	storage.NewGauge("HeapInuse", float64(metr.HeapInuse))
	storage.NewGauge("HeapObjects", float64(metr.HeapObjects))
	storage.NewGauge("HeapReleased", float64(metr.HeapReleased))
	storage.NewGauge("HeapSys", float64(metr.HeapSys))
	storage.NewGauge("LastGC", float64(metr.LastGC))
	storage.NewGauge("Lookups", float64(metr.Lookups))
	storage.NewGauge("MCacheInuse", float64(metr.MCacheInuse))
	storage.NewGauge("MCacheSys", float64(metr.MCacheSys))
	storage.NewGauge("MSpanInuse", float64(metr.MSpanInuse))
	storage.NewGauge("MSpanSys", float64(metr.MSpanSys))
	storage.NewGauge("Mallocs", float64(metr.Mallocs))
	storage.NewGauge("NextGC", float64(metr.NextGC))
	storage.NewGauge("NumForcedGC", float64(metr.NumForcedGC))
	storage.NewGauge("NumGC", float64(metr.NumGC))
	storage.NewGauge("OtherSys", float64(metr.OtherSys))
	storage.NewGauge("PauseTotalNs", float64(metr.PauseTotalNs))
	storage.NewGauge("StackInuse", float64(metr.StackInuse))
	storage.NewGauge("StackSys", float64(metr.StackSys))
	storage.NewGauge("Sys", float64(metr.Sys))
	storage.NewGauge("TotalAlloc", float64(metr.TotalAlloc))
	storage.NewGauge("RandomValue", RandomValue)
	//}
}

func CreateIntMetric(storage storage.MemStorage) {
	//return map[string]int64{
	storage.NewCounter("PollCount", PollCount)
	//}
}
