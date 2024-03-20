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

func SendMetricJSON(sugar zap.SugaredLogger, res []model.Metrics) {
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
	for key, val := range f {
		value := val
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "gauge", Value: &value})
		SendMetricJSON(sugar, modelmetrics)
	}
	for key, val := range i {
		value := val
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "counter", Delta: &value})
		SendMetricJSON(sugar, modelmetrics)
	}
	return nil
}
