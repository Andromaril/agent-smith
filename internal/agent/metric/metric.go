package metric

import (
	"encoding/json"
	"fmt"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// func SendGaugeMetric(name string, value float64) {
// 	client := resty.New()
// 	url := fmt.Sprintf("http://%s/update/gauge/%s/%v", flag.FlagRunAddr, name, value)
// 	//fmt.Print(url)
// 	_, err := client.R().Post(url)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func SendMetricJSON(sugar zap.SugaredLogger, res model.Metrics, client *resty.Client, url string) {
	//client := resty.New()
	// //jsonData, err := json.Marshal(res)

	// if err != nil {
	// 	panic(err)
	// }
	//url := fmt.Sprintf("http://%s/update/", flag.FlagRunAddr)
	//fmt.Print(url)
	jsonData, err := json.Marshal(res)
	if err != nil {
		sugar.Errorw("JSON marshaling failed", err)
	}

	if err != nil {
		panic(err)
	}
	_, err1 := client.R().SetHeader("Content-Type", "application/json").SetBody(jsonData).Post(url)
	if err1 != nil {
		panic(err)
	}
}

// func SendCounterMetric(name string, value int64) {
// 	client := resty.New()
// 	url := fmt.Sprintf("http://%s/update/counter/%s/%v", flag.FlagRunAddr, name, value)
// 	_, err := client.R().Post(url)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func SendAllMetric() error {
// 	f := creator.CreateFloatMetric()
// 	i := creator.CreateIntMetric()
// 	for key, value := range f {
// 		SendGaugeMetric(key, value)
// 	}
// 	for key, value := range i {
// 		SendCounterMetric(key, value)
// 	}
// 	return nil
// }

// func SendAllMetricJSON() error {
// 	f := creator.CreateFloatMetric()
// 	i := creator.CreateIntMetric()
// 	for key, value := range f {
// 		resp := model.Metrics{
// 			ID:    key,
// 			MType: "gauge",
// 			Value: &value,
// 		}
// 		SendMetricJSON(&resp)

// 	}
// 	for key, value := range i {
// 		resp := model.Metrics{
// 			ID:    key,
// 			MType: "counter",
// 			Delta: &value,
// 		}
// 		SendMetricJSON(&resp)
// 	}
// 	return nil
// }

func SendAllMetricJSON2(sugar zap.SugaredLogger, client *resty.Client) error {
	f := creator.CreateFloatMetric()
	i := creator.CreateIntMetric()
	//client := resty.New()
	url := fmt.Sprintf("http://%s/update/", flag.FlagRunAddr)
	for key, value := range f {
		resp := model.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value,
		}
		// jsonData, err := json.Marshal(resp)

		// if err != nil {
		// 	panic(err)
		// }
		SendMetricJSON(sugar, resp, client, url)
		// _, err1 := client.R().SetBody(jsonData).SetHeader("Content-Type", "application/json").Post(url)
		// if err1 != nil {
		// 	panic(err1)
		// }
	}
	for key, value := range i {
		resp := model.Metrics{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
		// jsonData, err := json.Marshal(resp)

		// if err != nil {
		// 	panic(err)
		// }
		// _, err1 := client.R().SetBody(jsonData).SetHeader("Content-Type", "application/json").Post(url)
		// if err1 != nil {
		// 	panic(err1)
		// }
		SendMetricJSON(sugar, resp, client, url)
	}
	return nil
}
