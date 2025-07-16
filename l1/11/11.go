package main

import "fmt"

/*1
Пересечение множеств
Реализовать пересечение двух неупорядоченных множеств (например, двух слайсов) — т.е. вывести элементы, присутствующие и в первом, и во втором.

Пример:
A = {1,2,3}
B = {2,3,4}
Пересечение = {2,3}
*/

// если принять множество набором уникальных элементов, можно сделать значение мапы bool или struct{}{} для экономии места и проверять просто наличие ключа в мапе
func main() {
	A := []int{1, 2, 3}
	B := []int{2, 3, 4}
	numbers := make(map[int]int)
	result := make([]int, 0, min(len(A), len(B)))
	for _, val := range A {
		numbers[val] += 1
	}

	for _, val := range B {
		count, ok := numbers[val]
		if ok && count > 0 {
			result = append(result, val)
			numbers[val] -= 1
		}

	}
	fmt.Println(result)
}
