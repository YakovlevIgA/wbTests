package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

// Выведется error, потому что внутри err помимо значения nil, есть ссылка на тип customError. Так как одно поле не nil, то сравнение err с nil даст false.
