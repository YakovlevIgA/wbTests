package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

/*

Программа должна корректно завершаться по нажатию Ctrl+C (SIGINT).

Выберите и обоснуйте способ завершения работы всех горутин-воркеров при получении сигнала прерывания.

Подсказка: можно использовать контекст (context.Context) или канал для оповещения о завершении.

*/

// todo: Обоснование - контекст быстрее реагирует и безопаснее, чем управляющий канал.
func workk(ctx context.Context, workerId int, in chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %v: context cancelled, exiting\n", workerId)
			return
		case val, ok := <-in:
			if !ok {
				fmt.Printf("Worker %v: channel closed, exiting\n", workerId)
				return
			}
			fmt.Printf("Worker %v said: %v\n", workerId, val)
		}
	}
}

func main() {
	const numWorkers = 3
	jobs := make(chan string)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Ловим ctrl+c
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		fmt.Println("Получен сигнал прерывания, отмена контекста")
		cancel()
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workk(ctx, i, jobs, wg)
	}

	go func() {
		for i := 0; i < 500; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Writer: context cancelled, stopping writing")
				close(jobs)
				return
			default:
				jobs <- strconv.Itoa(i)
				time.Sleep(100 * time.Millisecond) // имитация работы
			}
		}
		close(jobs)
	}()

	wg.Wait()

}
