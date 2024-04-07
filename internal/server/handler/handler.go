package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/go-chi/chi/v5"
)

func GetMetricJSON(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.MType == "counter" {
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
		if r.MType == "gauge" {
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

func GaugeandCounterJSON(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		res.Header().Set("Content-Type", "application/json")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&r); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.MType == "counter" {
			err := m.NewCounter(r.ID, *r.Delta)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			value, err := m.GetCounter(r.ID)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
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
		if r.MType == "gauge" {
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

func GaugeandCounter(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pattern := chi.URLParam(req, "pattern")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")
		res.Header().Set("Content-Type", "text/plain")
		if pattern == "counter" {
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				m.NewCounter(name, v)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else if pattern == "gauge" {
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				m.NewGauge(name, v)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}
	}
}

func GetMetric(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pattern := chi.URLParam(req, "pattern")
		name := chi.URLParam(req, "name")
		res.Header().Set("Content-Type", "text/plain")
		if pattern == "counter" {
			r, err := m.GetCounter(name)
			if err != nil {
				http.Error(res, "Incorrect metrics", http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
		} else if pattern == "gauge" {
			r, err := m.GetGauge(name)
			if err != nil {
				http.Error(res, "Incorrect metrics", http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}

	}
}

func GetHTMLMetric(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		//s := m.PrintMetric()
		gauge, err := m.GetFloatMetric()
		if err != nil {
			http.Error(res, "Incorrect metrics", http.StatusNotFound)
			return
		}
		counter, err2 := m.GetIntMetric()
		if err2 != nil {
			http.Error(res, "Incorrect metrics", http.StatusNotFound)
			return
		}
		var result string
		for k1, v1 := range gauge {
			result += fmt.Sprintf("%s: %v\n", k1, v1)
		}
		for k2, v2 := range counter {
			result += fmt.Sprintf("%s: %v\n", k2, v2)
		}
		tem := "<html> <head> <title> Metric page</title> </head> <body> <h1> List of metrics </h1> <p>" + html.EscapeString(result) + "</p> </body> </html>"
		res.Header().Set("Content-Type", "text/html")
		res.Write([]byte(tem))
	}
}

func Ping(db storagedb.Interface) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := db.Ping()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	}
}

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
			if models.MType == "gauge" {
				gauge = append(gauge, model.Gauge{Key: models.ID, Value: *models.Value})
			} else if models.MType == "counter" {
				counter = append(counter, model.Counter{Key: models.ID, Value: *models.Delta})
			}
		}
		err2 := db.CounterAndGaugeUpdateMetrics(gauge, counter)
		if err2 != nil {
			res.WriteHeader(http.StatusBadRequest)
		} else {
			res.WriteHeader(http.StatusOK)
		}
	}
}
