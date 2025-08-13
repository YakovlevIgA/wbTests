package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	// Параметры командной строки
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--timeout=<duration>] host port\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	// Подключение с таймаутом
	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, "Подключено к %s\n", address)

	// Канал для сигнализации о завершении
	done := make(chan struct{})

	// Чтение из соединения → STDOUT
	go func() {
		defer close(done)
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && !isClosedError(err) {
			fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		}
	}()

	// Чтение из STDIN → соединение
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := fmt.Fprintln(conn, scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка записи: %v\n", err)
				return
			}
		}
		// EOF (Ctrl+D) → закрыть соединение на запись
		conn.(*net.TCPConn).CloseWrite()
	}()

	// Ждём закрытия соединения
	<-done
	fmt.Fprintln(os.Stderr, "Соединение закрыто")
}

// isClosedError проверяет, что ошибка связана с закрытием соединения
func isClosedError(err error) bool {
	if err == io.EOF {
		return true
	}
	if ne, ok := err.(net.Error); ok && !ne.Timeout() {
		return true
	}
	return false
}
