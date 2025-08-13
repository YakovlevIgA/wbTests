package main

/*
Объяснить порядок вывода
*/

import "fmt"

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())        // 2
	fmt.Println(anotherTest()) // 1
}

// В случае с test() вывод будет 2, так как x - именованная возвращаемая переменная и при return не происходит фиксации значения.
// В anotherTest() происходит фиксация значения на return и вызванный дефер на вывод не влияет, поэтому получается 1.
