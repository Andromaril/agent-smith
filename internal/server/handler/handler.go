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

func ListMetric(r *http.Request) (string, string, *float64, *int64, error) {
	contentType := r.Header.Get("content-type")
	var types string
	var name string
	var value2 *float64
	var delta *int64
	//fmt.Println(contentType)
	if contentType == "application/json" {
		var req model.Metrics
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return "", "", nil, nil, err
		}
		types = req.MType
		name = req.ID
		value2 = req.Value
		delta = req.Delta
	} else {
		types = chi.URLParam(r, "pattern")
		name = chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")
		if types == "counter" {
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				delta = &v
				value2 = nil
			}
		} else if types == "gauge" {
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				value2 = &v
				delta = nil
			}
		}
	}
	return types, name, value2, delta, nil
}

func GaugeandCounter(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		types, name, value, delta, err := ListMetric(req)
		if err != nil {
			panic(err)
		}
		contentType := req.Header.Get("Content-Type")
		if types == "counter" {
			if delta != nil {
				m.NewCounter(name, *delta)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else if types == "gauge" {
			if value != nil {
				m.NewGauge(name, *value)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}
		if contentType == "application/json" {
			resp := model.Metrics{
				ID:    name,
				MType: types,
				Delta: delta,
				Value: value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}
}

func GetMetric(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		types, name2, value, delta, err := ListMetric(req)
		if err != nil {
			panic(err)
		}
		contentType := req.Header.Get("Content-Type")
		//fmt.Println(contentType)
		if types == "counter" {
			r, err := m.GetCounter(name2)
			if err != nil {
				http.Error(res, "Incorrect metrics", http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
		} else if types == "gauge" {
			r, err := m.GetGauge(name2)
			if err != nil {
				http.Error(res, "Incorrect metrics", http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}

		if contentType == "application/json" {
			resp := model.Metrics{
				ID:    name2,
				MType: types,
				Delta: delta,
				Value: value,
			}
			enc := json.NewEncoder(res)
			if err := enc.Encode(resp); err != nil {
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}
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
