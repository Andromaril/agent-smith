package flag

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr    string
	ReportInterval int64
	PollInterval   int64
	KeyHash        string
)

func ParseFlags() {
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&KeyHash, "k", "", "key HashSHA256")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		n, err := strconv.ParseInt(envReportInterval, 10, 64)
		if err != nil {
			panic(err)
		}
		ReportInterval = n
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		n, err := strconv.ParseInt(envPollInterval, 10, 64)
		if err != nil {
			panic(err)
		}
		PollInterval = n
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		KeyHash = envKey
	}
}
