package main

import (
	"fmt"
	"sync"
)

/*
Реализовать структуру-счётчик, которая будет инкрементироваться в конкурентной среде (т.е. из нескольких горутин).
По завершению программы структура должна выводить итоговое значение счётчика.
Подсказка: вам понадобится механизм синхронизации, например, sync.Mutex или sync/Atomic для безопасного инкремента.
*/

type counter struct {
	quantity int
	mu       sync.Mutex
	wg       sync.WaitGroup
}

func (c *counter) count() {
	defer c.wg.Done()
	c.mu.Lock()
	c.quantity += 1
	c.mu.Unlock()
}

func main() {
	c := counter{}

	for i := 0; i < 100; i++ {
		c.wg.Add(1)
		go c.count()

	}

	c.wg.Wait()
	fmt.Println(c.quantity)
}
