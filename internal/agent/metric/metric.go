package metric

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func SendMetricJSON(sugar zap.SugaredLogger, res *[]model.Metrics) {
	jsonData, err := json.Marshal(res)
	if err != nil {
		sugar.Errorw("marshalling error", err)
	}
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	zb.Write(jsonData)
	zb.Close()
	client := resty.New()
	url := fmt.Sprintf("http://%s/updates/", flag.FlagRunAddr)
	client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post(url)
}

func SendAllMetricJSON(sugar zap.SugaredLogger, storage storage.MemStorage) error {
	f, _ := storage.GetFloatMetric()
	i, _ := storage.GetIntMetric()
	modelmetrics := make([]model.Metrics, 0)
	for key, value := range f {
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "gauge", Value: &value})
		// resp := model.Metrics{
		// 	ID:    key,
		// 	MType: "gauge",
		// 	Value: &value,
		// }
		//modelmetrics = append(modelmetrics, resp)
		SendMetricJSON(sugar, &modelmetrics)
	}
	for key, value := range i {
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "gauge", Delta: &value})
		// resp := model.Metrics{
		// 	ID:    key,
		// 	MType: "counter",
		// 	Delta: &value,
		// }
		//modelmetrics = append(modelmetrics, resp)

		SendMetricJSON(sugar, &modelmetrics)
	}
	//SendMetricJSON(sugar, modelmetrics)
	return nil
	// client := resty.New()
	// gauges,_ := storage.GetFloatMetric()
	// counters,_ := storage.GetIntMetric()
	// metrics := make([]model.Metrics, 0)
	// url := fmt.Sprintf("http://%s/updates/", flag.FlagRunAddr)
	// buf := bytes.NewBuffer(nil)
	// zb := gzip.NewWriter(buf)
	// for key, val := range gauges {
	// 	value := val
	// 	metrics = append(metrics, model.Metrics{ID: key, MType: "gauge", Value: &value})
	// }
	// for key, val := range counters {
	// 	value := val
	// 	metrics = append(metrics, model.Metrics{ID: key, MType: "counter", Delta: &value})
	// }
	// jsonMetric, err := json.Marshal(metrics)
	// if err != nil {
	// 	return err
	// }
	// _, err = zb.Write(jsonMetric)
	// if err != nil {
	// 	return err
	// }
	// if err = zb.Close(); err != nil {
	// 	return err
	// }
	// r := client.NewRequest()
	// r.Header.Set("Content-Encoding", "gzip")
	// r.SetBody(buf)
	// if _, err = r.Post(url); err != nil {
	// 	return err
	// }
	// return nil
}
