// Package serverflag считывает флаги сервера
package serverflag

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr     string // адрес запуска сервиса
	StoreInterval   int64  // время между сохранениями метрик в файл
	FileStoragePath string // путь для файла с метриками
	Restore         bool   // определяет загружать или нет метрики с файла
	Databaseflag    string // адрес бд
	KeyHash         string // хеш
	CryptoKey       string // приватный ключ
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "path name")
	flag.Int64Var(&StoreInterval, "i", 300, "interval to save to disk")
	flag.BoolVar(&Restore, "r", true, "download files")
	flag.StringVar(&Databaseflag, "d", "", "database path")
	flag.StringVar(&KeyHash, "k", "", "key HashSHA256")
	flag.StringVar(&CryptoKey, "crypto-key", "", "key private")
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
}
