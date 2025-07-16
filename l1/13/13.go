package main

import "fmt"

/*
Поменять местами два числа без использования временной переменной.

Подсказка: примените сложение/вычитание или XOR-обмен.
*/
func main() {
	a, b := 5, 10
	c, d := 5, 10
	a = a + b // v.1
	b = a - b
	a = a - b

	c, d = d, c // v.2

	fmt.Println(a, b)
	fmt.Println(c, d)
}
