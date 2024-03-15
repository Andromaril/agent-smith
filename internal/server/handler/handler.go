package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func GetMetricJSON(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		//res.Header().Set("Content-Type", "application/json")
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
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		//w.Write(resp)
	}
}

func GaugeandCounterJSON(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var r model.Metrics
		//res.Header().Set("Content-Type", "application/json")
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
			//res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
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
			//res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
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
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
	}
}

func GaugeandCounter(m *storage.MemStorage) http.HandlerFunc {
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

func GetMetric(m *storage.MemStorage) http.HandlerFunc {
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

func GetHTMLMetric(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		s := m.PrintMetric()
		tem := "<html> <head> <title> Metric page</title> </head> <body> <h1> List of metrics </h1> <p>" + html.EscapeString(s) + "</p> </body> </html>"
		res.Header().Set("Content-Type", "text/html")
		res.Write([]byte(tem))
	}
}
