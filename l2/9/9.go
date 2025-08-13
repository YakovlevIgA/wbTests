package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
Написать функцию Go, осуществляющую примитивную распаковку строки, содержащей повторяющиеся символы/руны.

Примеры работы функции:

Вход: "a4bc2d5e"
Выход: "aaaabccddddde"

Вход: "abcd"
Выход: "abcd" (нет цифр — ничего не меняется)

Вход: "45"
Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)

Вход: ""
Выход: "" (пустая строка -> пустая строка)

Дополнительное задание
Поддерживать escape-последовательности вида \:

Вход: "qwe\4\5"
Выход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)

Вход: "qwe\45"
Выход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)

Требования к реализации
Функция должна корректно обрабатывать ошибочные случаи (возвращать ошибку, например, через error), и проходить unit-тесты.

Код должен быть статически анализируем (vet, golint).
*/

func unpack(s string) string {
	runes := []rune(s)
	var builder strings.Builder
	var digits string
	var lastLetterPosition int
	var haveLetters bool

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			i++
			builder.WriteRune(runes[i])
			lastLetterPosition = i
			continue
		}
		if unicode.IsLetter(runes[i]) == true { // обработка буквы
			builder.WriteRune(runes[i])
			lastLetterPosition = i
			haveLetters = true
		}

		if unicode.IsDigit(runes[i]) == true { // обработка цифры
			digits += string(runes[i])
			if i < len(runes)-1 {
				if unicode.IsDigit(runes[i+1]) == true {
					continue
				}
			}

			num, err := strconv.Atoi(digits)
			if err != nil {
				fmt.Println(err)
			}
			if haveLetters {
				for j := 0; j < num-1; j++ {
					builder.WriteRune(runes[lastLetterPosition])
				}
				digits = ""
			}
		}
	}

	return builder.String()
}

func main() {
	a := "a4bc2d5e"
	b := "abcd"
	c := "45"
	d := ""
	e := "a10bc4e"
	f := "1"
	g := `qwe\4\5`
	h := `qwe\45`
	fmt.Println(unpack(a))
	fmt.Println(unpack(b))
	fmt.Println(unpack(c))
	fmt.Println(unpack(d))
	fmt.Println(unpack(e))
	fmt.Println(unpack(f))
	fmt.Println(unpack(g))
	fmt.Println(unpack(h))

}
