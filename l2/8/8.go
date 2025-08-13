package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	time, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.Println("Ошибка получения времени через NTP:", err)
		os.Exit(1)
	}

	fmt.Println("Точное время (NTP):", time)
}
