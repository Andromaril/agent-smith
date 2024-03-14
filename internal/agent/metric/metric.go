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
)

func SendMetricJSON(res *model.Metrics) {
	jsonData, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	zb.Write(jsonData)
	zb.Close()
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flag.FlagRunAddr)
	client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(buf).Post(url)
}

func SendAllMetricJSON(storage storage.MemStorage) error {
	f := storage.GetFloatMetric()
	i := storage.GetIntMetric()

	for key, value := range f {
		resp := model.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}
		SendMetricJSON(&resp)
	}
	for key, value := range i {
		resp := model.Metrics{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}

		SendMetricJSON(&resp)
	}
	return nil
}
