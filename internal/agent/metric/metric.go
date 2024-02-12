package metric

import (
	"fmt"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/go-resty/resty/v2"
)

func SendGaugeMetric(name string, value float64) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/gauge/%s/%v", flag.FlagRunAddr, name, value)
	//fmt.Print(url)
	_, err := client.R().Post(url)
	if err != nil {
		panic(err)
	}
}

func SendCounterMetric(name string, value int64) {
	client := resty.New()
	url := fmt.Sprintf("http://%s/update/counter/%s/%v", flag.FlagRunAddr, name, value)
	_, err := client.R().Post(url)
	if err != nil {
		panic(err)
	}
}

func SendAllMetric() error {
	f := creator.CreateFloatMetric()
	i := creator.CreateIntMetric()
	for key, value := range f {
		SendGaugeMetric(key, value)
	}
	for key, value := range i {
		SendCounterMetric(key, value)
	}
	return nil
}
