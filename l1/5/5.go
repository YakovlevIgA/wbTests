package main

import (
	"context"
	"fmt"
	"time"
)

/*
Разработать программу, которая будет последовательно отправлять значения в канал, а с другой стороны канала – читать эти значения.
По истечении N секунд программа должна завершаться.
*/

func main() {
	ch := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go func() {
		defer close(ch)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				i++
			}
		}
	}()

	for val := range ch {
		fmt.Println(val)

	}

}
