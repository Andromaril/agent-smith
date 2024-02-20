package handler

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func GaugeandCounter(m *storage.MemStorage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pattern := chi.URLParam(req, "pattern")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")
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
