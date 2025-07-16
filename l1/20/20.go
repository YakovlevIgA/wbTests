package main

import "fmt"

/*
Разработать программу, которая переворачивает порядок слов в строке.

Пример: входная строка:

«snow dog sun», выход: «sun dog snow».

Считайте, что слова разделяются одиночным пробелом. Постарайтесь не использовать дополнительные срезы, а выполнять операцию «на месте».
*/

func reverse2(str string) string {
	runes := []rune(str)

	for i := 0; i < len(runes)/2; i++ {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}

	left := 0
	for right := 0; right <= len(runes); right++ {
		if right == len(runes) || runes[right] == ' ' {
			for l, r := left, right-1; l < r; l, r = l+1, r-1 {
				runes[l], runes[r] = runes[r], runes[l]
			}
			left = right + 1
		}
	}

	return string(runes)
}

func main() {
	someString := "snow dog sun"
	fmt.Println(reverse2(someString))
}
