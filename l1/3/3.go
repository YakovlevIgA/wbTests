package main

import (
	"fmt"
	"strconv"
	"sync"
)

/*
Реализовать постоянную запись данных в канал (в главной горутине).

Реализовать набор из N воркеров, которые читают данные из этого канала и выводят их в stdout.

Программа должна принимать параметром количество воркеров и при старте создавать указанное число горутин-воркеров.


*/

func work(workerId int, in chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range in {
		fmt.Printf("Worker %v said: %v\n", workerId, val)
	}

}

func main() {
	const numWorkers = 3
	jobs := make(chan string)
	wg := &sync.WaitGroup{}

	go func() {
		for i := 0; i < 500; i++ {
			jobs <- strconv.Itoa(i)
		}
		close(jobs)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go work(i, jobs, wg)
	}
	wg.Wait()

}
