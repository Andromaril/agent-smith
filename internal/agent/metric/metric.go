// Package metric содержит функцию, отправлющую метрики в json-формате
package metric

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/andromaril/agent-smith/internal/constant"
	"github.com/andromaril/agent-smith/internal/errormetric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	pb "github.com/andromaril/agent-smith/pkg/proto"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Metric struct {
	MetricClient pb.MetricClient
}

func (m *Metric) AddMetricClient(mc pb.MetricClient) {
	m.MetricClient = mc
}

// SendMetricJSON функция, отправляющая метрики в формате json по эндпоинту /updates/
func SendMetricJSON(sugar zap.SugaredLogger, res []model.Metrics) error {
	conn, err := net.Dial("tcp", flag.FlagRunAddr)
	if err != nil {
		sugar.Errorw(
			"error when send mentric",
			"error", err,
		)
	}
	localAddress := conn.LocalAddr().(*net.TCPAddr)
	reqHeader := localAddress.IP.To4().String()
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
			SetHeader("X-Real-IP", reqHeader).
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
		pub, err := x509.ParsePKIXPublicKey(pemDecode.Bytes)
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error decode public key %w", e)
		}
		buf2, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), buf.Bytes())
		if err != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error decode public key %w", e)
		}
		_, err2 := client.R().SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("X-Real-IP", reqHeader).
			SetBody(buf2).
			Post(url)
		if err2 != nil {
			e := errormetric.NewMetricError(err)
			return fmt.Errorf("error send request %w", e)
		}
	}
	_, err2 := client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", reqHeader).
		SetBody(buf).
		Post(url)
	if err2 != nil {
		e := errormetric.NewMetricError(err)
		return fmt.Errorf("error send request %w", e)
	}
	return nil
}

// SendMetricGRPC функция, отправляющая метрики grpc по эндпоинту /updates/
func SendMetricGRPC(m Metric, sugar zap.SugaredLogger, res []model.Metrics) error {
	if m.MetricClient == nil {
		return errors.New("metricClient has a nil pointer")
	}
	gauge := make([]*pb.Gauge, 0)
	counter := make([]*pb.Counter, 0)
	for _, models := range res {
		if models.MType == constant.Gauge {
			gauge = append(gauge, &pb.Gauge{Key: models.ID, Value: *models.Value})

		} else if models.MType == constant.Counter {
			counter = append(counter, &pb.Counter{Key: models.ID, Value: *models.Delta})
		}
	}
	_, err := m.MetricClient.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{
		Gauge:   gauge,
		Counter: counter,
	})
	sugar.Infow(
		"update",
		"Gauge", gauge, "Counter", counter)
	return err
}
