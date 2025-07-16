package main

import (
	"fmt"
	"sync"
)

/*
Реализовать безопасную для конкуренции запись данных в структуру map.

Подсказка: необходимость использования синхронизации (например, sync.Mutex или встроенная concurrent-map).

Проверьте работу кода на гонки (util go run -race).
*/

type memory struct {
	mappy map[int]string
	mu    sync.Mutex
}

func newMemory() *memory {
	return &memory{
		mappy: make(map[int]string),
	}
}

func (m *memory) insert(key int, value string, wg *sync.WaitGroup) {
	defer wg.Done()
	m.mu.Lock()
	m.mappy[key] = value
	m.mu.Unlock()

}

func main() {
	wg := &sync.WaitGroup{}
	mem := newMemory()
	value := "yes"
	for key := 0; key < 100; key++ {
		wg.Add(1)
		go mem.insert(key, value, wg)

	}
	wg.Wait()
	fmt.Println(mem.mappy)
}
