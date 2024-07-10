// Package serverflag считывает флаги сервера
package serverflag

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	FlagRunAddr     string `json:"address"`
	Restore         bool   `json:"restore"`
	StoreInterval   int64  `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	Databaseflag    string `json:"database_dsn"`
	CryptoKey       string `json:"crypto_key"`
	ConfigKey       string
	TrustedSubnet   string `json:"trusted_subnet"`
	GrpcKey         string `json:"grpc_key"`
}

var (
	FlagRunAddr     string // адрес запуска сервиса
	StoreInterval   int64  // время между сохранениями метрик в файл
	FileStoragePath string // путь для файла с метриками
	Restore         bool   // определяет загружать или нет метрики с файла
	Databaseflag    string // адрес бд
	KeyHash         string // хеш
	CryptoKey       string // приватный ключ
	ConfigKey       string // файл с конфигом в формате json
	TrustedSubnet   string // строковое представление бесклассовой адресации (CIDR)
	GrpcKey         string // адрес запуска сервиса grpc
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path name")
	flag.Int64Var(&StoreInterval, "i", 300, "interval to save to disk")
	flag.BoolVar(&Restore, "r", true, "download files")
	flag.StringVar(&Databaseflag, "d", "", "database path")
	flag.StringVar(&KeyHash, "k", "", "key HashSHA256")
	flag.StringVar(&CryptoKey, "crypto-key", "", "key private")
	flag.StringVar(&ConfigKey, "c", "", "json-file flag")
	flag.StringVar(&TrustedSubnet, "t", "", "CIDR")
	flag.StringVar(&GrpcKey, "g", "", "GRPC")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}
	if envDatabaseflag := os.Getenv("DATABASE_DSN"); envDatabaseflag != "" {
		Databaseflag = envDatabaseflag
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
	if envKey := os.Getenv("KEY"); envKey != "" {
		KeyHash = envKey
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		CryptoKey = envCryptoKey
	}
	if envConfigKey := os.Getenv("CONFIG"); envConfigKey != "" {
		ConfigKey = envConfigKey
	}
	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		TrustedSubnet = envTrustedSubnet
	}
	if envGrpcKey := os.Getenv("GRPC"); envGrpcKey != "" {
		GrpcKey = envGrpcKey
	}

	if ConfigKey != "" {
		c, err := os.ReadFile(ConfigKey)
		if err != nil {
			panic(err)
		}
		var conf Config
		err = json.Unmarshal(c, &conf)
		if err != nil {
			log.Fatal(err)
		}
		if FlagRunAddr == "localhost:8080" {
			FlagRunAddr = conf.FlagRunAddr
		}
		if Restore {
			Restore = conf.Restore
		}
		if StoreInterval == 300 {
			StoreInterval = conf.StoreInterval
		}
		if FileStoragePath == "" {
			FileStoragePath = conf.FileStoragePath
		}
		if Databaseflag == "" {
			Databaseflag = conf.Databaseflag
		}
		if CryptoKey == "" {
			CryptoKey = conf.CryptoKey
		}
		if TrustedSubnet == "" {
			TrustedSubnet = conf.TrustedSubnet
		}
	}
}
