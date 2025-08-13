package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseFields(fieldStr string) map[int]bool {
	fields := make(map[int]bool)
	parts := strings.Split(fieldStr, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			limits := strings.Split(part, "-")
			if len(limits) != 2 {
				continue
			}
			start, err1 := strconv.Atoi(limits[0])
			end, err2 := strconv.Atoi(limits[1])
			if err1 != nil || err2 != nil || start > end {
				continue
			}
			for i := start; i <= end; i++ {
				fields[i-1] = true // -1: индексация с 0
			}
		} else {
			i, err := strconv.Atoi(part)
			if err != nil {
				continue
			}
			fields[i-1] = true
		}
	}
	return fields
}

func main() {
	fieldStr := flag.String("f", "", "fields (e.g., 1,3-5)")
	delimiter := flag.String("d", "\t", "delimiter (default: tab)")
	separatedOnly := flag.Bool("s", false, "only lines with delimiter")

	flag.Parse()

	if *fieldStr == "" {
		fmt.Fprintln(os.Stderr, "Error: -f is required")
		os.Exit(1)
	}

	fields := parseFields(*fieldStr)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		if *separatedOnly && !strings.Contains(line, *delimiter) {
			continue
		}

		parts := strings.Split(line, *delimiter)
		var selected []string

		for i, col := range parts {
			if fields[i] {
				selected = append(selected, col)
			}
		}

		fmt.Println(strings.Join(selected, *delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}
}
