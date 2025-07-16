package main

import (
	"fmt"
	"math"
)

/*
Разработать программу нахождения расстояния между двумя точками на плоскости.
Точки представлены в виде структуры Point с инкапсулированными (приватными) полями x, y (типа float64) и конструктором.
Расстояние рассчитывается по формуле между координатами двух точек.
Подсказка: используйте функцию-конструктор NewPoint(x, y), Point и метод Distance(other Point) float64.
*/
type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func (p *Point) Distance(o *Point) float64 {
	dx := o.x - p.x
	dy := o.y - p.y
	return math.Sqrt(dx*dx + dy*dy)
}

func main() {
	a := NewPoint(2, 2)
	b := NewPoint(3, 3)
	fmt.Println(a.Distance(b))
}
