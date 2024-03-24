package router

import (
	"net/http"

	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/middleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Router(newMetric storage.Storage, sugar zap.SugaredLogger) {
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	r.Use(logging.WithLogging(sugar))
	r.Route("/value", func(r chi.Router) {
		r.Post("/", handler.GetMetricJSON(newMetric))
		r.Get("/{pattern}/{name}", handler.GetMetric(newMetric))
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handler.GaugeandCounterJSON(newMetric))
		r.Post("/{pattern}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	r.Get("/", handler.GetHTMLMetric(newMetric))
	r.Get("/ping", handler.Ping(newMetric.(storagedb.Interface)))
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", handler.Update(newMetric))
	})
	if err := http.ListenAndServe(serverflag.FlagRunAddr, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")

	}
}
