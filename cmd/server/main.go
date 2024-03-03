package main

import (
	"net/http"

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
	//fmt.Printf(flag.FileStoragePath)
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
	newMetric := storage.NewMemStorage()
	//var i int64
	//if serverflag.StoreInterval != 0 {
	//for i = 0; ; i++ {
	//time.Sleep(time.Second)
	//if i%serverflag.StoreInterval == 0 {
	//}
	//}
	//storage.Save(newMetric)
	//storage.Load(newMetric)
	//} else {
	//newMetric.Save(serverflag.FileStoragePath)
	//fmt.Printf(flag.FileStoragePath)
	//}
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
	// var i int64
	// if serverflag.StoreInterval != 0 {
	// 	for i = 0; ; i++ {
	// 		if i%serverflag.StoreInterval == 0 {
	// 			newMetric.Save(serverflag.FileStoragePath)
	// 			fmt.Printf(serverflag.FileStoragePath)
	// 			time.Sleep(time.Second * time.Duration(serverflag.StoreInterval))
	// 		} else {
	// 			newMetric.Save(serverflag.FileStoragePath)
	// 		}
	// 	}
	// }
	//for {
	go newMetric.Save(serverflag.FileStoragePath, serverflag.StoreInterval)
	//fmt.Print(123)
	//time.Sleep(time.Second * time.Duration(serverflag.StoreInterval))
	//}

	if err := http.ListenAndServe(serverflag.FlagRunAddr, r); err != nil {
		sugar.Fatalw(err.Error(), "event", "start server")

	}

}

// func saveMetrics(path string, interval int) {
// 	for {
// 		storage.Save(path)
// 		time.Sleep(time.Duration(interval) * time.Second)
// 	}
// }
