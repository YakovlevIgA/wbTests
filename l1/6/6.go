package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
Реализовать все возможные способы остановки выполнения горутины.

Классические подходы: выход по условию, через канал уведомления, через контекст, прекращение работы runtime.Goexit() и др.

Продемонстрируйте каждый способ в отдельном фрагменте кода.
*/
var stop bool

func withCondition(wg *sync.WaitGroup) {
	defer wg.Done()
	res := 0
	for {
		if stop {
			fmt.Printf("goroutine stopped because of the condition with res = %v\n", res)
			return
		}
		res += 1
	}
}

func withChannel(ch chan struct{}, wg *sync.WaitGroup) {
	res := 0
	for {
		select {
		case <-ch:
			fmt.Printf("goroutine stopped because by sending signal to the channel with res = %v\n", res)
			wg.Done()
			return
		default:
			res += 1
		}
	}

}
func withContext(ctx context.Context, wg *sync.WaitGroup) {
	res := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("goroutine stopped by context with res = %v\n", res)
			wg.Done()
			return
		default:
			res += 1
		}
	}
}

func withGoexit(wg *sync.WaitGroup) {
	res := 0
	for i := 0; i < 100; i++ {
		res += 1
	}
	fmt.Printf("goroutine stopped by Goexit with res = %v\n", res)
	wg.Done()
	runtime.Goexit()
}

func withClosingCh(in chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	res := 0
	for val := range in {
		res += val
	}
	fmt.Printf("goroutine stopped by closing chan with res = %v\n", res)
}

func main() {
	wg := &sync.WaitGroup{}
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ с условием
	wg.Add(1)
	go withCondition(wg)
	time.Sleep(time.Millisecond * 10)
	stop = true
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ с каналом
	wg.Add(1)
	ch := make(chan struct{})
	go withChannel(ch, wg)
	time.Sleep(time.Millisecond * 10)
	ch <- struct{}{}
	close(ch)
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ с контекстом
	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	go withContext(ctx, wg)
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ с Goexit
	wg.Add(1)
	go withGoexit(wg)
	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ с закрытием канала
	wg.Add(1)
	ch2 := make(chan int)
	go withClosingCh(ch2, wg)
	for i := 1; ; {
		if i > 100 {
			close(ch2)
			break
		}
		ch2 <- i
		i++
	}

	//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ завершить main

	wg.Wait()

}
