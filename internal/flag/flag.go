package flag

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	FlagRunAddr    string
	ReportInterval int64
	PollInterval   int64
)

func ParseFlags() {
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		n, err := strconv.ParseInt(envReportInterval, 10, 64)
		ReportInterval = n
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(ReportInterval)
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		n, err := strconv.ParseInt(envPollInterval, 10, 64)
		PollInterval = n
		if err != nil {
			panic(err)
		}
	}
}
