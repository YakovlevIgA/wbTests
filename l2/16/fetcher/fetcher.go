package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Fetch(u *url.URL, client *http.Client) ([]byte, string, error) {
	start := time.Now()
	logURL := u.String()
	fmt.Printf("[FETCHING] %s\n", logURL)

	req, err := http.NewRequest("GET", logURL, nil)
	if err != nil {
		fmt.Printf("[ERROR] Creating request: %v\n", err)
		return nil, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		if isTimeout(err) {
			fmt.Printf("[TIMEOUT] %s (%v)\n", logURL, time.Since(start))
		} else {
			fmt.Printf("[ERROR] %s: %v\n", logURL, err)
		}
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Reading body %s: %v\n", logURL, err)
		return nil, "", err
	}

	contentType := resp.Header.Get("Content-Type")
	fmt.Printf("[DONE] %s (%v, %s)\n", logURL, time.Since(start), contentType)

	return body, contentType, nil
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if ue, ok := err.(interface{ Timeout() bool }); ok {
		return ue.Timeout()
	}
	return false
}
