package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/andromaril/agent-smith/internal/agent/creator"
	"github.com/andromaril/agent-smith/internal/agent/metric"
	"github.com/andromaril/agent-smith/internal/flag"
	"github.com/andromaril/agent-smith/internal/model"
	"github.com/andromaril/agent-smith/internal/retry"
	pb "github.com/andromaril/agent-smith/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var sugar zap.SugaredLogger
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func worker(m *metric.Metric, wg *sync.WaitGroup, jobs <-chan []model.Metrics, sugar zap.SugaredLogger) {
	defer wg.Done()
	for j := range jobs {
		operation := func() error {
			if flag.GrpcKey != "" {
				err := metric.SendMetricGRPC(*m, sugar, j)
				return err
			}
			err := metric.SendMetricJSON(sugar, j)
			return err
		}
		err2 := retry.Retry(operation)

		if err2 != nil {
			sugar.Errorw(
				"error when send mentric",
				"error", err2)
		}
		time.Sleep(time.Second * time.Duration(flag.ReportInterval))
	}
}

func main() {
	flag.ParseFlags()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()
	var m *metric.Metric = new(metric.Metric)

	if flag.GrpcKey != "" {
		conn, err := grpc.NewClient(flag.GrpcKey, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			sugar.Errorw(
				"error grpc agent",
				"error", err,
			)
		}
		defer conn.Close()
		// получаем переменную интерфейсного типа UsersClient,
		// через которую будем отправлять сообщения
		c := pb.NewMetricClient(conn)
		m.AddMetricClient(c)
		if m.MetricClient == nil {
			sugar.Errorw(
				"error grpc agent, nil metricclient",
			)
		}
		//c1.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{})
		sugar.Infow(
			"start gRPC agent",
		)
	}

	sugar.Infow(
		"Starting agent")
	sugar.Infow(
		"Starting agent",
		"Build version:", buildVersion, "Build date:", buildDate, "Build commit:", buildCommit)
	var wg sync.WaitGroup
	ratelimit := flag.RateLimit
	jobs := make(chan []model.Metrics, runtime.GOMAXPROCS(0))
	wg.Add(ratelimit)
	go creator.CreateFloatMetric(jobs)
	go creator.AddNewMetric(jobs)
	defer close(jobs)
	for w := 1; w <= ratelimit; w++ {
		go worker(m, &wg, jobs, sugar)
	}
	wg.Wait()
	defer wg.Done()
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		err := metric.SendMetricJSON(sugar, <-jobs)
		if err != nil {
			sugar.Errorw(
				"error when send mentric",
				"error", err,
			)
		}
		close(idleConnsClosed)
	}()
	<-idleConnsClosed
	fmt.Println("Agent stop")
}
