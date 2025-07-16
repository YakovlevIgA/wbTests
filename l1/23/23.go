package main

import "fmt"

/*
Удалить i-ый элемент из слайса. Продемонстрируйте корректное удаление без утечки памяти.

Подсказка: можно сдвинуть хвост слайса на место удаляемого элемента (copy(slice[i:], slice[i+1:])) и уменьшить длину слайса на 1.
*/

func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	slice = slice[:len(slice)-1]
	return slice
}

func main() {
	s := []int{10, 20, 30, 40, 50}
	fmt.Println("Исходный слайс:", s)

	i := 2

	s = remove(s, i)
	fmt.Println("После удаления:", s)
}
