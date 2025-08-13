/*
Что выведет программа?

Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.
*/
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)        // nil, как значение переменной
	fmt.Println(err == nil) // false, так как у err остается ссылка на тип - т.е. одно из полей не nil.

}
