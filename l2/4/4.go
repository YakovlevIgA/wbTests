package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i <= 10; i++ {
			ch <- i
		}
	}()
	for n := range ch {
		println(n)
	}
}

// от 0 до 9, затем дедлок, так как range ждет закрытия, которого не происходит.
