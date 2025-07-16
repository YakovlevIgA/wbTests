package main

import "fmt"

/*Рассмотреть следующий код и ответить на вопросы: к каким негативным последствиям он может привести и как это исправить?

Приведите корректный пример реализации.

var justString string

func someFunc() {
	v := createHugeString(1 &lt;&lt; 10)
	justString = v[:100]
}

func main() {
	someFunc()
}
Вопрос: что происходит с переменной justString?
*/



// todo: Из-за глобальной области видимости juststring будет срезом v и после выхода из функции, но в первой реализации GC не освободит переменную v,
// todo: так как на неё еще будут ссылаться и получится утечка памяти. Лучше явно скопировать в juststring строку через указанный синтаксис или copy.

var justString string

func someFunc() {
	v := createHugeString(1 &lt;&lt; 10)
	justString = string([]rune(v[:100]))
}

func main() {
	someFunc()
}