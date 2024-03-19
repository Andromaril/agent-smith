package serverflag

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr     string
	StoreInterval   int64
	FileStoragePath string
	Restore         bool
)

func ParseFlags() {
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
