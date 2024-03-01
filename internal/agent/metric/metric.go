package metric

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/go-resty/resty/v2"
)

func SendMetricJSON(res *model.Metrics) {
	jsonData, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	// var b bytes.Buffer
	// w, err := flate.NewWriter(&b, flate.BestCompression)
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = w.Write(jsonData)

	// if err != nil {
	// 	panic(err)
	// }
	// err = w.Close()
	// if err != nil {
	//    panic(err)
	// }

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, gzErr := gz.Write(jsonData); gzErr != nil {
		return
	}
	if gzErr := gz.Close(); gzErr != nil {
		return
	}
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flag.FlagRunAddr)
	//fmt.Print(url)
	client.R().SetHeader("Content-Type", "application/json").SetBody(b.Bytes()).Post(url)
	// if err1 != nil {
	// 	panic(err1)
	// }
}

func SendAllMetricJSON2() error {
	f := creator.CreateFloatMetric()
	i := creator.CreateIntMetric()

	for key, value := range f {
		resp := model.Metrics{
			ID:    key,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		SendMetricJSON(&resp)
	}
	for key, value := range i {
		resp := model.Metrics{
			ID:    key,
			MType: "counter",
			Delta: &value,
			Value: nil,
		}

		SendMetricJSON(&resp)
	}
	return nil
}
