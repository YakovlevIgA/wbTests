package main

import (
	"fmt"
	"strings"
)

/*
Разработать программу, которая проверяет, что все символы в строке встречаются один раз (т.е. строка состоит из уникальных символов).

Вывод: true, если все символы уникальны, false, если есть повторения. Проверка должна быть регистронезависимой, т.е. символы в разных регистрах считать одинаковыми.

Например: "abcd" -> true, "abCdefAaf" -> false (повторяются a/A), "aabcd" -> false.

Подумайте, какой структурой данных удобно воспользоваться для проверки условия.
*/

func check(s string) bool {
	lib := make(map[rune]struct{})
	runes := []rune(strings.ToLower(s))
	for _, val := range runes {
		_, ok := lib[val]
		if !ok {
			lib[val] = struct{}{}
		} else {
			return false
		}

	}
	return true
}

func main() {
	str1 := "abcd"
	str2 := "aAbcd"
	fmt.Println(check(str1))
	fmt.Println(check(str2))
}
