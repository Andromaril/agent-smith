package flag

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr     string
	ReportInterval  int64
	PollInterval    int64
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
)

func ParseFlags() {
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path name")
	flag.Int64Var(&StoreInterval, "i", 300, "interval to save to disk")
	flag.BoolVar(&Restore, "r", true, "download files")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
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
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		n, err := strconv.ParseInt(envStoreInterval, 10, 64)
		if err != nil {
			panic(err)
		}
		StoreInterval = n
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		n, err := strconv.ParseBool(envRestore)
		if err != nil {
			panic(err)
		}
		Restore = n
	}
}
