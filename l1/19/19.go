package main

import "fmt"

/*
Разработать программу, которая переворачивает подаваемую на вход строку.

Например: при вводе строки «главрыба» вывод должен быть «абырвалг».

Учтите, что символы могут быть в Unicode (русские буквы, emoji и пр.), то есть просто iterating по байтам может не подойти — нужен срез рун ([]rune).
*/

func reverse(str string) string {
	runes := []rune(str)
	for i := 0; i < len(runes)/2; i++ {
		runes[i], runes[len(runes)-i-1] = runes[len(runes)-i-1], runes[i]
	}

	return string(runes)
}

func main() {
	someString := "главрыба"
	fmt.Println(reverse(someString))
}
