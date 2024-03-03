package main

import (
	"net/http"
	"time"

	"github.com/andromaril/agent-smith/internal/flag"
	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/middleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	flag.ParseFlags()
	logger, err1 := zap.NewDevelopment()
	if err1 != nil {
		panic(err1)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow(
		"Starting server",
		"addr", flag.FlagRunAddr,
	)
	newMetric := storage.NewMemStorage()
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
	if err := http.ListenAndServe(flag.FlagRunAddr, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")
	}
	storage.RestoreData(newMetric)
	var i int64
	if flag.StoreInterval != 0 {
		for i = 0; ; i++ {
			time.Sleep(time.Second)
			if i%flag.StoreInterval == 0 {
				storage.Save(newMetric)
				time.Sleep(time.Second * time.Duration(flag.StoreInterval))
			}
		}
		//storage.Save(newMetric)
		//storage.Load(newMetric)
	}
}
