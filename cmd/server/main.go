package main

import (
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	newMetric := storage.NewMemStorage()
	r := chi.NewRouter()
	r.Route("/value", func(r chi.Router) {
		r.Get("/{change}/{name}", handler.GetMetric(newMetric))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{change}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	r.Get("/", handler.GetHtmlMetric(newMetric))
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
