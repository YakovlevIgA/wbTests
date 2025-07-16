package main

import (
	"fmt"
	"math/big"
)

/*
Разработать программу, которая перемножает, делит, складывает, вычитает две числовых переменных a, b, значения которых > 2^20 (больше 1 миллион).
Комментарий: в Go тип int справится с такими числами, но обратите внимание на возможное переполнение для ещё больших значений.
Для очень больших чисел можно использовать math/big.
*/

func main() {

	a := big.NewInt(0)
	b := big.NewInt(0)

	a.SetString("15000000000000000000000000", 10) // число очень большое
	b.SetString("20000000000000000000000000", 10)

	// Сумма
	sum := big.NewInt(0).Add(a, b)
	// Разность
	diff := big.NewInt(0).Sub(a, b)
	// Произведение
	mul := big.NewInt(0).Mul(a, b)
	// Частное
	div := big.NewInt(0).Div(a, b)

	fmt.Printf("a = %s\nb = %s\n", a.String(), b.String())
	fmt.Printf("Сумма: %s\n", sum.String())
	fmt.Printf("Разность: %s\n", diff.String())
	fmt.Printf("Произведение: %s\n", mul.String())
	fmt.Printf("Частное: %s\n", div.String())
}
