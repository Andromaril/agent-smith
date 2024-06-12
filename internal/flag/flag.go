// Package flag считывает флаги агента
package flag

import (
	"flag"
	"os"
	"strconv"
)

var (
	FlagRunAddr    string // адрес запуска агента
	ReportInterval int64  // время между отправкой метрик
	PollInterval   int64  // время сбора метрик
	KeyHash        string // хеш
	RateLimit      int    // количество горутин
	CryptoKey      string // публичный ключ
)

// ParseFlags для флагов либо переменных окружения
func ParseFlags() {
	flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
	flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&KeyHash, "k", "", "key HashSHA256")
	flag.IntVar(&RateLimit, "l", 2, "rate limit")
	flag.StringVar(&CryptoKey, "crypto-key", "", "key public")
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
}
