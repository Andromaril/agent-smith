package handler

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/utils"
	"github.com/go-chi/chi/v5"
)


func GaugeandCounter(m *storage.MemStorage) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		spath := utils.ParseURL(req.URL)
		change := chi.URLParam(req, "change")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")
		if change == "counter" {
			if value1, err := strconv.ParseInt(value, 10, 64); err == nil {
				m.NewCounter(name, value1)
			} else {
				http.Error(res, "Incorrect metrics" , http.StatusBadRequest)
			}
		} else if change == "gauge" {
			if value1, err := strconv.ParseFloat(value, 64); err == nil {
				m.NewGauge(spath[3], value1)
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
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		change := chi.URLParam(req, "change")
		name := chi.URLParam(req, "name")
		if change == "counter" {
			r, err := m.GetCounter(name)
			if err != nil {
				http.Error(res, "Incorrect metrics" , http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
			} else if change == "gauge" {
				r, err := m.GetGauge(name)
				if err != nil {
					http.Error(res, "Incorrect metrics" , http.StatusNotFound)
					return
				}
				res.Write([]byte(fmt.Sprint(r)))
			} else {
				http.Error(res, "Incorrect metrics" , http.StatusBadRequest)
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