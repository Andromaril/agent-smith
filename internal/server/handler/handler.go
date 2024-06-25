// Package handler содержит хэндлеры
package handler

import (
	"fmt"
	"html"
	"net/http"
	"strconv"

	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/go-chi/chi/v5"
	"github.com/andromaril/agent-smith/internal/constant"
)

// GaugeandCounter позволяет создать новую метрику по post запросу update/{pattern}/{name}/{value} принимая данные в url-параметрах запроса
func GaugeandCounter(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pattern := chi.URLParam(req, "pattern")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")
		res.Header().Set("Content-Type", "text/plain")
		if pattern == constant.Counter {
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				m.NewCounter(name, v)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else if pattern == constant.Gauge{
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

// GetMetric по GET запросу к value/{pattern}/{name} выводит значение метрики
func GetMetric(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		pattern := chi.URLParam(req, "pattern")
		name := chi.URLParam(req, "name")
		res.Header().Set("Content-Type", "text/plain")
		if pattern == constant.Counter {
			r, err := m.GetCounter(name)
			if err != nil {
				http.Error(res, "Incorrect metrics", http.StatusNotFound)
				return
			}
			res.Write([]byte(fmt.Sprint(r)))
		} else if pattern == constant.Gauge {
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

// GetHTMLMetric по GET запросу к /value/ выводит значение всех метрик в отдельной html-странице
func GetHTMLMetric(m storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
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

// Ping по GET-запросу к /ping проверяет состояние базы данных
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
