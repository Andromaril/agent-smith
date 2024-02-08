package flag

import (
    "flag"
)
var (
	FlagRunAddr string
    ReportInterval int64
    PollInterval int64
)
// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера

// parseFlags обрабатывает аргументы командной строки 
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
    // регистрируем переменную flagRunAddr 
    // как аргумент -a со значением :8080 по умолчанию
    flag.Int64Var(&ReportInterval, "r", 10, "time to sleep for report interval")
    flag.Int64Var(&PollInterval, "p", 2, "time to sleep for poll interval")
    flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()
}
