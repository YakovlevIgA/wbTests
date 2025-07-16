package main

import (
	"fmt"
	"time"
)

/*
Реализовать собственную функцию sleep(duration) аналогично встроенной функции time.Sleep, которая приостанавливает выполнение текущей горутины.

Важно: в отличии от настоящей time.Sleep, ваша функция должна именно блокировать выполнение (например, через таймер или цикл), а не просто вызывать time.Sleep :) — это упражнение.

Можно использовать канал + горутину, или цикл на проверку времени (не лучший способ, но для обучения).
*/

func Sleep(d time.Duration) {
	ch := make(chan struct{})
	go func() {
		time.AfterFunc(d, func() {
			close(ch)
		})
	}()
	<-ch
}

func main() {
	fmt.Println("started")

	Sleep(1 * time.Second)
	fmt.Println("ended")
}
