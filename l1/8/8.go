package main

/*
Дана переменная типа int64. Разработать программу, которая устанавливает i-й бит этого числа в 1 или 0.

Пример: для числа 5 (0101₂) установка 1-го бита в 0 даст 4 (0100₂).

Подсказка: используйте битовые операции (|, &^).
*/

import (
	"fmt"
)

func setBitTo1(n int64, i uint) int64 {
	return n | (1 << i)
}

func setBitTo0(n int64, i uint) int64 {
	return n &^ (1 << i)
}

func main() {
	var n int64 = 5
	var i uint = 0

	fmt.Printf("Исходное число: %b (%d)\n", n, n)

	n = setBitTo0(n, i)
	fmt.Printf("После установки %d-го бита в 0: %b (%d)\n", i, n, n)

	n = setBitTo1(n, i)
	fmt.Printf("После установки %d-го бита в 1: %b (%d)\n", i, n, n)
}
