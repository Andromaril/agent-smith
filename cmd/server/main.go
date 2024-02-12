package main

import (
	"net/http"

	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	flag.ParseFlags()
	newMetric := storage.NewMemStorage()
	r := chi.NewRouter()
	r.Route("/value", func(r chi.Router) {
		r.Get("/{pattern}/{name}", handler.GetMetric(newMetric))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{pattern}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	r.Get("/", handler.GetHTMLMetric(newMetric))
	err := http.ListenAndServe(flag.FlagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
