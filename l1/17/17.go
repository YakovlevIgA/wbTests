package main

import "fmt"

/*
Реализовать алгоритм бинарного поиска встроенными методами языка.
Функция должна принимать отсортированный слайс и искомый элемент, возвращать индекс элемента или -1, если элемент не найден.
Подсказка: можно реализовать рекурсивно или итеративно, используя цикл for.
*/
func biSerch(sl []int, target int) int {
	left, right := 0, len(sl)-1

	for left <= right {
		mid := left + (right-left)/2

		if sl[mid] == target {
			return mid
		} else if sl[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1 // не найден

}

func main() {
	sl := []int{2, 4, 7, 10, 23, 60, 78, 444, 234235, 423423423}
	fmt.Println(biSerch(sl, 234235))

}
