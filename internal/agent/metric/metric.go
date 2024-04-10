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
	url := fmt.Sprintf("http://%s/updates/", flag.FlagRunAddr)
	if flag.KeyHash != "" {
		hash := hmac.New(sha256.New, []byte(flag.KeyHash))
		hash.Write(jsonData)
		dst := hex.EncodeToString(hash.Sum(jsonData))
		_, err2 := client.R().SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("HashSHA256", dst).
			SetBody(buf).
			Post(url)
		if err2 != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error send request %q", e.Error())
		}
	}
	_, err2 := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post(url)
	if err2 != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error send request %q", e.Error())
	}
	return nil
}
