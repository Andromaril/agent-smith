package handler

import (
	//"fmt"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/utils"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)


func GaugeandCounter(m *storage.MemStorage) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		spath := utils.ParseURL(req.URL)
		if len(spath) != 5 {
			http.Error(res, "Not found metrics", http.StatusNotFound)
			return
		}
		change := chi.URLParam(req, "change")
		name := chi.URLParam(req, "name")
		value := chi.URLParam(req, "value")
		//if len(spath) == 5 {
		if change == "counter" {
			if value1, err := strconv.ParseInt(value, 10, 64); err == nil {
				//value := chi.URLParam(req, spath[4])
				m.NewCounter(name, value1)
				//fmt.Println(value1)
				//fmt.Println(m)
			} else {
				http.Error(res, "Incorrect metrics" , http.StatusBadRequest)
			}
		} else if change == "gauge" {
			if value1, err := strconv.ParseFloat(value, 64); err == nil {
				m.NewGauge(spath[3], value1)
				//fmt.Println(value1)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}
	
	//} 
	//else {
		//http.Error(res, "Not found metrics", http.StatusNotFound)
		//return
	//}	
}
}