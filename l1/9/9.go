package main

import "fmt"

/*
Разработать конвейер чисел. Даны два канала: в первый пишутся числа x из массива, во второй – результат операции x*2.
После этого данные из второго канала должны выводиться в stdout.
То есть, организуйте конвейер из двух этапов с горутинами: генерация чисел и их обработка.
Убедитесь, что чтение из второго канала корректно завершается.
*/

func main() {
	first := make(chan int)
	second := make(chan int)
	x := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	go func() {
		for _, val := range x {
			first <- val
		}
		close(first)
	}()

	go func() {
		for val := range first {
			second <- val * 2
		}
		close(second)
	}()

	for res := range second {
		fmt.Println(res)
	}

}
