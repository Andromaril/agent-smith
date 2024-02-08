package main

import (
	"fmt"
	"net/http"
	//"runtime/metrics"

	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/go-chi/chi/v5"
)

func main() {
	flag.ParseFlags()
	newMetric := storage.NewMemStorage()
	r := chi.NewRouter()
	r.Route("/value", func(r chi.Router) {
		r.Get("/{change}/{name}", handler.GetMetric(newMetric))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{change}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	r.Get("/", handler.GetHTMLMetric(newMetric))
	fmt.Println("Running server on", flag.FlagRunAddr)
	err := http.ListenAndServe(flag.FlagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
