package serverflag

// import (
// 	"flag"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	"github.com/caarlos0/env"
// )

// type ServerConfig struct {
// 	Address         string
// 	StoreInterval   time.Duration
// 	FileStoragePath string
// 	Restore         bool
// }

// type ServerEnvConfig struct {
// 	Address         string `env:"ADDRESS"`
// 	StoreInterval   string `env:"STORE_INTERVAL"`
// 	FileStoragePath string `env:"FILE_STORAGE_PATH"`
// 	Restore         string `env:"RESTORE"`
// }

// var (
// 	flagRunAddr         string
// 	flagLogLevel        string
// 	flagStoreInterval   int
// 	flagFileStoragePath string
// 	flagRestore         bool
// )

// func Flags() (*ServerConfig, error) {
// 	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
// 	flag.StringVar(&flagLogLevel, "l", "info", "log level")
// 	flag.IntVar(&flagStoreInterval, "i", 300, "STORE INTERVAL")
// 	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
// 	flag.BoolVar(&flagRestore, "r", true, "restore")

// 	// парсим переданные серверу аргументы в зарегистрированные переменные
// 	flag.Parse()

// 	var config ServerEnvConfig

// 	err := env.Parse(&config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resultConfig := ServerConfig{}

// 	if config.Address == "" {
// 		resultConfig.Address = flagRunAddr
// 	} else {
// 		resultConfig.Address = config.Address
// 	}

// 	if config.StoreInterval == "" {
// 		resultConfig.StoreInterval = time.Second * time.Duration(flagStoreInterval)
// 	} else {
// 		v, err := strconv.ParseInt(config.StoreInterval, 10, 32)
// 		if err != nil {
// 			panic(fmt.Errorf("problem setup config.StoreInterval %w", err))
// 		}
// 		resultConfig.StoreInterval = time.Second * time.Duration(v)
// 	}

// 	if config.FileStoragePath == "" {
// 		resultConfig.FileStoragePath = flagFileStoragePath
// 	} else {
// 		resultConfig.FileStoragePath = config.FileStoragePath
// 	}

// 	if config.Restore == "" {
// 		resultConfig.Restore = flagRestore
// 	} else {
// 		restore, err := strconv.ParseBool(config.Restore)
// 		if err != nil {
// 			panic(err)
// 		}
// 		resultConfig.Restore = restore
// 	}

// 	return &resultConfig, nil
// }

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr     string
	//ReportInterval  int64
	//PollInterval    int64
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
