package main

import (
	"net/http"
	"time"

	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/middleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	serverflag.ParseFlags()
	logger, err1 := zap.NewDevelopment()
	if err1 != nil {
		panic(err1)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	sugar.Infow(
		"Starting server",
		"addr", serverflag.FlagRunAddr,
	)
	newMetric := storage.NewMemStorage(serverflag.StoreInterval == 0, serverflag.FileStoragePath)

	if serverflag.Restore {
		newMetric.Load(serverflag.FileStoragePath)
	}
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
	//В задании указано, что StoreInterval = 0 делает запись синхронной.
	//Запись в бесконечном цикле - не то же самое.
	if serverflag.StoreInterval != 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(serverflag.StoreInterval))
			newMetric.Save(serverflag.FileStoragePath)
		}()
	}

	if err := http.ListenAndServe(serverflag.FlagRunAddr, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")

	}
}
