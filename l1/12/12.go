package main

import (
	"fmt"
)

/*
Имеется последовательность строк: ("cat", "cat", "dog", "cat", "tree"). Создать для неё собственное множество.

Ожидается: получить набор уникальных слов. Для примера, множество = {"cat", "dog", "tree"}.
*/

func makeUnique(in []string) {
	book := make(map[string]struct{})
	res := make([]string, 0, len(in))
	for _, word := range in {
		_, ok := book[word]
		if !ok {
			book[word] = struct{}{}
			res = append(res, word)
		}
	}
	fmt.Println(res)
}

func main() {
	input := []string{"cat", "cat", "dog", "cat", "tree"}
	makeUnique(input)
}
