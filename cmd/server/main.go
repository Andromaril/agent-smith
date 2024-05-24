// Package main запускает сервис
package main

import (
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "net/http/pprof"

	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/midleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/start"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

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
	db, newMetric := start.Start()
	if serverflag.Restore {
		newMetric.Load(serverflag.FileStoragePath)
	}
	defer db.Close()
	r := chi.NewRouter()
	r.Use(midleware.GzipMiddleware)
	if serverflag.KeyHash != "" {
		r.Use(midleware.HashMiddleware(serverflag.KeyHash))
	}
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
	r.Mount("/debug", middleware.Profiler())
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
