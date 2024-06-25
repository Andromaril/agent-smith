package handlerdb

import (
	"encoding/json"
	"net/http"

	"github.com/andromaril/agent-smith/internal/constant"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
)

// GetMetricJSON 1 метрику в json-формате по запросу к /value
func GetMetricJSON(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.MType == constant.Counter {
			value, err := m.GetCounter(r.ID)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			}
			resp := model.Metrics{
				ID:    r.ID,
				MType: r.MType,
				Delta: &value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				return
			}
		}
		if r.MType == constant.Gauge {
			value, err := m.GetGauge(r.ID)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
			}
			resp := model.Metrics{
				ID:    r.ID,
				MType: r.MType,
				Value: &value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				return
			}
		}
		res.WriteHeader(http.StatusOK)
	}
}

// GaugeandCounterJSON позволяет создать новую метрику по post запросу /update/ принимая данные в json-формате
func GaugeandCounterJSON(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.MType == constant.Counter {
			err := m.NewCounter(r.ID, *r.Delta)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			value, err := m.GetCounter(r.ID)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			resp := model.Metrics{
				ID:    r.ID,
				MType: r.MType,
				Delta: &value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				return
			}
		}
		if r.MType == constant.Gauge {
			err := m.NewGauge(r.ID, *r.Value)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			value, err := m.GetGauge(r.ID)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				return
			}
			resp := model.Metrics{
				ID:    r.ID,
				MType: r.MType,
				Value: &value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				return
			}
		}
		res.WriteHeader(http.StatusOK)
	}
}

// Update по Post запросу /updates обновляет метрики в базе данных
func Update(db storagedb.Interface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r := make([]model.Metrics, 0)
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		gauge := make([]model.Gauge, 0)
		counter := make([]model.Counter, 0)
		for _, models := range r {
			if models.MType == constant.Gauge {
				gauge = append(gauge, model.Gauge{Key: models.ID, Value: *models.Value})
			} else if models.MType == constant.Counter {
				counter = append(counter, model.Counter{Key: models.ID, Value: *models.Delta})
			}
		}
		err2 := db.CounterAndGaugeUpdateMetrics(gauge, counter)
		if err2 != nil {
			res.WriteHeader(http.StatusBadRequest)
		}
		res.WriteHeader(http.StatusOK)
	}
}
