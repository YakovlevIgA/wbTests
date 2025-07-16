package main

import "fmt"

/*
Разработать программу, которая в runtime способна определить тип переменной, переданной в неё (на вход подаётся interface{}).

	Типы, которые нужно распознавать: int, string, bool, chan (канал).

Подсказка: оператор типа switch v.(type) поможет в решении.
*/
func checkType(param interface{}) {
	switch param.(type) {
	case int:
		fmt.Println("It's int")
	case string:
		fmt.Println("It's string")
	case bool:
		fmt.Println("It's bool")
	case chan int:
		fmt.Println("It's chan")

	}

}

func main() {
	var chislo int
	var stroka string
	var flag bool
	var ch chan int

	checkType(chislo)
	checkType(stroka)
	checkType(flag)
	checkType(ch)

}
