package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	logging "github.com/andromaril/agent-smith/internal/loger"
	"github.com/andromaril/agent-smith/internal/middleware"
	"github.com/andromaril/agent-smith/internal/server/handler"
	"github.com/andromaril/agent-smith/internal/server/storage"
	"github.com/andromaril/agent-smith/internal/server/storage/storagedb"
	"github.com/andromaril/agent-smith/internal/serverflag"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func main() {
	serverflag.ParseFlags()
	//fmt.Print(serverflag.Databaseflag)
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
	var newMetric storage.Storage
	var err error
	var db *sql.DB
	if serverflag.Databaseflag != "" {
		newMetric = &storagedb.StorageDB{Path: serverflag.Databaseflag}
		db, err = newMetric.Init(serverflag.Databaseflag, context.Background())
		// db, err = sql.Open("pgx", serverflag.Databaseflag)
		if err != nil {
			panic(err)
		}
		//defer db.Close()
	} else {
		newMetric = &storage.MemStorage{Gauge: map[string]float64{}, Counter: map[string]int64{}, WriteSync: serverflag.StoreInterval == 0, Path: serverflag.FileStoragePath}
	}
	// db, err := sql.Open("pgx", serverflag.Databaseflag)
	// if err != nil {
	// 	panic(err)
	// }
	defer db.Close()
	//newMetric.Init(serverflag.FileStoragePath, context.Background())
	// var newMetric *storage.MemStorage
	// if serverflag.Databaseflag != "" {
	// 	newMetric = storage.NewMemStorage(serverflag.StoreInterval == 0, serverflag.FileStoragePath)
	// } else {
	// 	newMetric = &storagedb.StorageDB{Path: serverflag.FileStoragePath}
	// }

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
	r.Get("/ping", handler.Ping(db))
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
