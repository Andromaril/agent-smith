package handler

import (
	//"fmt"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/utils"
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
		if spath[2] == "counter" {
			if counter2, err := strconv.ParseInt(spath[4], 10, 64); err == nil {
				m.NewCounter(spath[3], counter2)
				//fmt.Println(counter2)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else if spath[2] == "gauge" {
			if gauge2, err := strconv.ParseFloat(spath[4], 64); err == nil {
				m.NewGauge(spath[3], gauge2)
				//fmt.Println(gauge2)
			} else {
				http.Error(res, "Incorrect metrics", http.StatusBadRequest)
			}
		} else {
			http.Error(res, "Incorrect metrics", http.StatusBadRequest)
		}
	
	}
}