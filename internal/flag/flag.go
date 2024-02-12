package flag

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

var (
	FlagRunAddr    string
	ReportInterval int64
	PollInterval   int64
)

type Config struct {
	ReportInterval int64 `env:"REPORT_INTERVAL"`
	PollInterval   int64 `env:"POLL_INTERVAL"`
}

func ParseFlags() {
	var cfg Config
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

}
