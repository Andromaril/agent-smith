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
	//r.Route("/update", func(r chi.Router) {
        //r.Post("/", handler.GaugeandCounter(newMetric))
       // r.Route("/", func(r chi.Router) {
            //r.Post("/gauge", handler.GaugeandCounter(newMetric))
			//r.Route("/gauge", func(r chi.Router) {
				//r.Post("/gauge/", handler.GaugeandCounter(newMetric))
			//}        // GET /cars/renault
            //r.Post("/counter", handler.GaugeandCounter(newMetric)) // GET /cars/renault/duster
       // })
    //})
	r.Route("/update", func(r chi.Router) {
		r.Post("/{change}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	//newMetric := storage.NewMemStorage()
	//mux := http.NewServeMux()
	//mux.HandleFunc("/update/", handler.GaugeandCounter(newMetric))
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
	//log.Fatal(http.ListenAndServe(":8080", r))
}
