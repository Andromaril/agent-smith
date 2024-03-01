package metric

import (
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
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flag.FlagRunAddr)
	//fmt.Print(url)
	client.R().SetHeader("Content-Type", "application/json").SetBody(jsonData).Post(url)
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
