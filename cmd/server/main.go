// Package main запускает сервис
package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "net/http/pprof"

	"github.com/andromaril/agent-smith/internal/errormetric"
	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/middleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/handlerdb"
	"github.com/andromaril/agent-smith/internal/server/start"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
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
	sugar.Infow(
		"Starting server",
		"Build version:", buildVersion, "Build date:", buildDate, "Build commit:", buildCommit)
	db, newMetric := start.Start()
	if serverflag.Restore {
		newMetric.Load(serverflag.FileStoragePath)
	}
	defer db.Close()
	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware)
	if serverflag.KeyHash != "" {
		r.Use(middleware.HashMiddleware(serverflag.KeyHash))
	}
	if serverflag.CryptoKey != "" {
		data, err := os.ReadFile(serverflag.CryptoKey)
		if err != nil {
			e := errormetric.NewMetricError(err)
			sugar.Errorw(
				"error read file",
				"error", e,
			)
		}
		pemDecode, _ := pem.Decode(data)
		priv, err := x509.ParsePKCS1PrivateKey(pemDecode.Bytes)
		if err != nil {
			e := errormetric.NewMetricError(err)
			sugar.Errorw(
				"error parse",
				"error", e,
			)
		}
		r.Use(middleware.CryptoMiddleware(priv))
	}
	if serverflag.ConfigKey != "" {
		r.Use(middleware.IPMiddleware(serverflag.TrustedSubnet))
	}
	r.Use(logging.WithLogging(sugar))
	r.Route("/value", func(r chi.Router) {
		r.Post("/", handlerdb.GetMetricJSON(newMetric))
		r.Get("/{pattern}/{name}", handler.GetMetric(newMetric))
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlerdb.GaugeandCounterJSON(newMetric))
		r.Post("/{pattern}/{name}/{value}", handler.GaugeandCounter(newMetric))
	})
	r.Get("/", handler.GetHTMLMetric(newMetric))
	r.Get("/ping", handler.Ping(newMetric.(storagedb.Interface)))
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", handlerdb.Update(newMetric))
	})
	//r.Mount("/debug", middleware.Profiler())
	if serverflag.StoreInterval != 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(serverflag.StoreInterval))
			newMetric.Save(serverflag.FileStoragePath)
		}()
	}
	var srv = http.Server{Addr: serverflag.FlagRunAddr, Handler: r}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			sugar.Errorw(
				"HTTP server Shutdown",
				"error", err,
			)
		}
		err := newMetric.Save(serverflag.FileStoragePath)
		if err != nil {
			sugar.Errorw(
				"error save to file",
				"error", err,
			)
		}
		close(idleConnsClosed)
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		sugar.Fatalw(err.Error(), "event", "start server")
	}

	<-idleConnsClosed
	sugar.Infow(
		"Server Shutdown gracefully",
	)
}
