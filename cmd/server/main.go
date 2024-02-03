package main

import (
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"net/http"
)


func main() {
	newMetric := storage.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.GaugeandCounter(newMetric))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
