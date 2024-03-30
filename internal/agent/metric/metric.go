package metric

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func SendMetricJSON(sugar zap.SugaredLogger, res []model.Metrics) error {
	jsonData, err := json.Marshal(res)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error %q", e.Error())
	}
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	zb.Write(jsonData)
	zb.Close()
	client := resty.New()
	r := client.NewRequest()
	url := fmt.Sprintf("http://%s/updates/", flag.FlagRunAddr)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Encoding", "gzip")
	if flag.KeyHash != "" {
		hash := hmac.New(sha256.New, []byte(flag.KeyHash))
		hash.Write(jsonData)
		dst := hex.EncodeToString(hash.Sum(nil))
		r.Header.Set("HashSHA256", dst)
	}
	r.SetBody(buf)
	_, err2 := r.Post(url)
	if err2 != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error send request %q", e.Error())
	}
	return nil
}

func SendAllMetricJSON(sugar zap.SugaredLogger, storage storage.MemStorage) error {
	f, _ := storage.GetFloatMetric()
	i, _ := storage.GetIntMetric()
	modelmetrics := make([]model.Metrics, 0)
	for key, val := range f {
		value := val
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "gauge", Value: &value})
	}
	for key, val := range i {
		value := val
		modelmetrics = append(modelmetrics, model.Metrics{ID: key, MType: "counter", Delta: &value})
	}
	err := SendMetricJSON(sugar, modelmetrics)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error %q", e.Error())
	}
	return nil
}
