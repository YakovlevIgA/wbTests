package main

import (
	"fmt"
	"sync"
)

/*
Написать программу, которая конкурентно рассчитает значения квадратов чисел, взятых из массива [2,4,6,8,10], и выведет результаты в stdout.

Подсказка: запусти несколько горутин, каждая из которых возводит число в квадрат.

*/

func main() {
	arr := [5]int{2, 4, 6, 8, 10}
	/*
	   	for i := 0; i < len(arr); i++ {

	   		go func() {
	   			fmt.Println(arr[i] * arr[i])
	   		}()
	   	}
	   	time.Sleep(2 * time.Second)
	   }
	*/

	wg := &sync.WaitGroup{}

	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(arr[i] * arr[i])
		}()
	}
	wg.Wait()

}
