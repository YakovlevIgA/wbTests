package main

import "fmt"

/* Дана структура Human (с произвольным набором полей и методов).

Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).

Подсказка: используйте композицию (embedded struct), чтобы Action имел все методы Human.
*/

type Human struct {
	name   string
	age    int
	height int
}

type Action struct {
	Human
}

func (h Human) walk() {
	fmt.Println("walking")
}
func (h Human) jump() {
	fmt.Println("jumping")
}
func (h Human) run() {
	fmt.Println("running")
}

func (a Action) think() {
	fmt.Println("thinking")
}

func main() {
	Alex := Human{
		name:   "Alex",
		age:    25,
		height: 185,
	}
	AlexRoutine := Action{
		Human: Alex,
	}
	AlexRoutine.walk()
	AlexRoutine.think()

}
