package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

type Config struct {
	after      int
	before     int
	countOnly  bool
	ignoreCase bool
	invert     bool
	fixed      bool
	showLineNo bool
	pattern    string
}

type Grep struct {
	cfg    Config
	lines  []string
	output map[int]bool // строки, которые нужно вывести
	marked map[int]bool // строки, уже отмеченные
}

func parseFlags() Config {
	A := flag.Int("A", 0, "Print N lines After match")
	B := flag.Int("B", 0, "Print N lines Before match")
	C := flag.Int("C", 0, "Print N lines Context around match")
	c := flag.Bool("c", false, "Count of matching lines only")
	i := flag.Bool("i", false, "Ignore case")
	v := flag.Bool("v", false, "Invert match")
	F := flag.Bool("F", false, "Fixed string (not regex)")
	n := flag.Bool("n", false, "Show line number")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: grep [options] pattern [file]")
		os.Exit(1)
	}
	pattern := args[0]

	// C overrides A and B
	if *C > 0 {
		*A = *C
		*B = *C
	}

	return Config{
		after:      *A,
		before:     *B,
		countOnly:  *c,
		ignoreCase: *i,
		invert:     *v,
		fixed:      *F,
		showLineNo: *n,
		pattern:    pattern,
	}
}

func (g *Grep) compile() (*regexp.Regexp, error) {
	p := g.cfg.pattern
	if g.cfg.fixed {
		p = regexp.QuoteMeta(p)
	}
	if g.cfg.ignoreCase {
		p = `(?i)` + p
	}
	return regexp.Compile(p)
}

func (g *Grep) match(line string, re *regexp.Regexp) bool {
	found := re.MatchString(line)
	if g.cfg.invert {
		return !found
	}
	return found
}

func (g *Grep) collectMatches(re *regexp.Regexp) int {
	count := 0
	for i, line := range g.lines {
		if g.match(line, re) {
			count++
			for j := i - g.cfg.before; j <= i+g.cfg.after; j++ {
				if j >= 0 && j < len(g.lines) {
					g.output[j] = true
				}
			}
		}
	}
	return count
}

func (g *Grep) printMatches(count int) {
	if g.cfg.countOnly {
		fmt.Println(count)
		return
	}

	lastPrinted := -1
	for i := 0; i < len(g.lines); i++ {
		if g.output[i] {
			if i > 0 && lastPrinted != i-1 && lastPrinted != -1 {
				fmt.Println("--")
			}
			prefix := ""
			if g.cfg.showLineNo {
				prefix = fmt.Sprintf("%d:", i+1)
			}
			fmt.Println(prefix + g.lines[i])
			lastPrinted = i
		}
	}
}

func readInput(path string) ([]string, error) {
	var reader io.Reader
	if path == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	cfg := parseFlags()
	var filename string
	if len(flag.Args()) > 1 {
		filename = flag.Args()[1]
	}

	lines, err := readInput(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	grep := &Grep{
		cfg:    cfg,
		lines:  lines,
		output: make(map[int]bool),
		marked: make(map[int]bool),
	}

	re, err := grep.compile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid pattern:", err)
		os.Exit(1)
	}

	count := grep.collectMatches(re)
	grep.printMatches(count)
}
