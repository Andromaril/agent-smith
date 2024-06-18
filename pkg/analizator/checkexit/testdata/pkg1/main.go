package main

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	mulfunc(5)
	res, _ := mulfunc(5)
	os.Exit(1) // want "calling os.Exit"
	fmt.Println(res)
}
