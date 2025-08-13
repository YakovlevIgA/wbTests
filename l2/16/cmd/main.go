package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"wget/crawler"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: mirror <URL> [depth]")
	}
	startURL := os.Args[1]
	depth := 2
	if len(os.Args) >= 3 {
		d, err := strconv.Atoi(os.Args[2])
		if err == nil {
			depth = d
		}
	}

	u, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	cr := crawler.NewCrawler(u, "output", depth, &http.Client{
		Timeout: 10 * time.Second,
	}, 8)

	cr.Run()
}
