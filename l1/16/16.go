package main

import "fmt"

/*
Реализовать алгоритм быстрой сортировки массива встроенными средствами языка. Можно использовать рекурсию.
Подсказка: напишите функцию quickSort([]int) []int которая сортирует срез целых чисел.
Для выбора опорного элемента можно взять середину или первый элемент.
*/

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	pivot := arr[0]
	var less []int
	var greater []int

	for _, v := range arr[1:] {
		if v <= pivot {
			less = append(less, v)
		} else {
			greater = append(greater, v)
		}
	}

	return append(append(quickSort(less), pivot), quickSort(greater)...)
}

func main() {
	arr := []int{88, 4, 2, 4, 61, 8, 3, 8, 3, 5, 0, 7, 2, 5}
	fmt.Println(quickSort(arr))

}
