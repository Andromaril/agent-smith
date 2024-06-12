// Package metric модержит функцию, отправлющую метрики в json-формате
package metric

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// SendMetricJSON функция, отправляющая метрики в формате json по эндпоинту /updates/
func SendMetricJSON(sugar zap.SugaredLogger, res []model.Metrics) error {
	jsonData, err := json.Marshal(res)
	if err != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error %w", e)
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
			return fmt.Errorf("error send request %w", e)
		}
	}
	if flag.CryptoKey != "" {
		data, err := os.ReadFile(flag.CryptoKey)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error read file %w", e)
		}
		pemDecode, _ := pem.Decode(data)
		pub, err := x509.ParsePKCS1PublicKey(pemDecode.Bytes)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error decode public key %w", e)
		}
		buf2, err := rsa.EncryptPKCS1v15(rand.Reader, pub, buf.Bytes())
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error decode public key %w", e)
		}
		_, err2 := client.R().SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(buf2).
			Post(url)
		if err2 != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error send request %w", e)
		}
	}
	_, err2 := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(buf).
		Post(url)
	if err2 != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error send request %w", e)
	}
	return nil
}
