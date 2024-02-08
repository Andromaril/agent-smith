package flag

import (
    "flag"
)
var (
	FlagRunAddr string
)
// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера

// parseFlags обрабатывает аргументы командной строки 
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
    // регистрируем переменную flagRunAddr 
    // как аргумент -a со значением :8080 по умолчанию
    flag.StringVar(&FlagRunAddr, "a", ":8080", "address and port to run server")
    // парсим переданные серверу аргументы в зарегистрированные переменные
    flag.Parse()
}
