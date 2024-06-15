// Package flag считывает флаги агента
package flag

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

type Config struct {
	FlagRunAddr    string `json:"address"`
	ReportInterval int64  `json:"report_interval"`
	PollInterval   int64  `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

var (
	FlagRunAddr    string // адрес запуска агента
	ReportInterval int64  // время между отправкой метрик
	PollInterval   int64  // время сбора метрик
	KeyHash        string // хеш
	RateLimit      int    // количество горутин
	CryptoKey      string // публичный ключ
	ConfigKey      string // файл с конфигом в формате json
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&KeyHash, "k", "", "key HashSHA256")
	flag.IntVar(&RateLimit, "l", 2, "rate limit")
	flag.StringVar(&CryptoKey, "crypto-key", "", "key public")
	flag.StringVar(&ConfigKey, "-c", "", "json-file flag")
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
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		n, err := strconv.ParseInt(envRateLimit, 10, 64)
		if err != nil {
			panic(err)
		}
		RateLimit = int(n)
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		CryptoKey = envCryptoKey
	}
	if envConfigKey := os.Getenv("CONFIG"); envConfigKey != "" {
		ConfigKey = envConfigKey
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
		if ReportInterval == 10 {
			ReportInterval = conf.ReportInterval
		}
		if PollInterval == 2 {
			PollInterval = conf.PollInterval
		}
		if CryptoKey == "" {
			CryptoKey = conf.CryptoKey
		}
	}
}
