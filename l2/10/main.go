package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Реализовать упрощённый аналог UNIX-утилиты sort (сортировка строк).

Программа должна читать строки (из файла или STDIN) и выводить их отсортированными.

Обязательные флаги (как в GNU sort):

-k N — сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию).
Например, «sort -k 2» отсортирует строки по второму столбцу каждой строки.

-n — сортировать по числовому значению (строки интерпретируются как числа).

-r — сортировать в обратном порядке (reverse).

-u — не выводить повторяющиеся строки (только уникальные).

Дополнительные флаги:

-M — сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат.

-b — игнорировать хвостовые пробелы (trailing blanks).

-c — проверить, отсортированы ли данные; если нет, вывести сообщение об этом.

-h — сортировать по числовому значению с учётом суффиксов (например, К = килобайт, М = мегабайт — человекочитаемые размеры).

Программа должна корректно обрабатывать комбинации флагов (например, -nr — числовая сортировка в обратном порядке, и т.д.).

Необходимо предусмотреть эффективную обработку больших файлов.

Код должен проходить все тесты, а также проверки go vet и golint (понимание, что требуются надлежащие комментарии, имена и структура программы).
*/

// monthOrder хранит порядковые номера месяцев для -M
var monthOrder = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

// humanSize преобразует строку с человекочитаемыми суффиксами K,M,G в float64
func humanSize(s string) (float64, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	multiplier := 1.0
	switch last := s[len(s)-1]; last {
	case 'K', 'k':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'M', 'm':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'G', 'g':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return num * multiplier, nil
}

func main() {
	// Определение флагов
	column := flag.Int("k", 0, "Сортировать по колонке N (по умолчанию вся строка)")
	numeric := flag.Bool("n", false, "Числовая сортировка")
	reverse := flag.Bool("r", false, "Обратная сортировка")
	unique := flag.Bool("u", false, "Уникальные строки")
	month := flag.Bool("M", false, "Сортировка по названию месяца (Jan, Feb...)")
	ignoreTrailingBlanks := flag.Bool("b", false, "Игнорировать хвостовые пробелы")
	checkSorted := flag.Bool("c", false, "Проверить, отсортированы ли данные")
	human := flag.Bool("h", false, "Сортировка с человекочитаемыми суффиксами K,M,G")
	flag.Parse()

	var reader io.Reader
	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	var lines []string
	seen := map[string]bool{}

	for scanner.Scan() {
		line := scanner.Text()
		if *ignoreTrailingBlanks {
			line = strings.TrimRight(line, " \t")
		}
		if *unique {
			if seen[line] {
				continue
			}
			seen[line] = true
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	// Проверка на отсортированность
	if *checkSorted {
		if isSorted(lines, getComparator(*column, *numeric, *month, *human, *reverse)) {
			fmt.Println("Data is sorted.")
		} else {
			fmt.Println("Data is NOT sorted.")
		}
		return
	}

	// Сортировка
	sort.SliceStable(lines, func(i, j int) bool {
		return getComparator(*column, *numeric, *month, *human, *reverse)(lines[i], lines[j])
	})

	// Вывод
	for _, line := range lines {
		fmt.Println(line)
	}
}

// getComparator возвращает функцию сравнения с учётом флагов
func getComparator(k int, numeric, month, human, reverse bool) func(string, string) bool {
	return func(a, b string) bool {
		keyA := extractKey(a, k)
		keyB := extractKey(b, k)

		var result int
		if month {
			keyA = strings.Title(strings.ToLower(keyA))
			keyB = strings.Title(strings.ToLower(keyB))
			ma, oka := monthOrder[keyA]
			mb, okb := monthOrder[keyB]
			if oka && okb {
				if ma < mb {
					result = -1
				} else if ma > mb {
					result = 1
				} else {
					result = 0
				}
			} else if oka && !okb {
				result = 1
			} else if !oka && okb {
				result = -1
			} else {
				result = strings.Compare(keyA, keyB)
			}

		} else if human {
			na, ea := humanSize(keyA)
			nb, eb := humanSize(keyB)
			if ea == nil && eb == nil {
				if na < nb {
					result = -1
				} else if na > nb {
					result = 1
				} else {
					result = 0
				}
			} else {
				result = strings.Compare(keyA, keyB)
			}
		} else if numeric {
			na, ea := strconv.ParseFloat(keyA, 64)
			nb, eb := strconv.ParseFloat(keyB, 64)
			if ea == nil && eb == nil {
				if na < nb {
					result = -1
				} else if na > nb {
					result = 1
				} else {
					result = 0
				}
			} else {
				result = strings.Compare(keyA, keyB)
			}
		} else {
			result = strings.Compare(keyA, keyB)
		}

		if reverse {
			return result > 0
		}
		return result < 0
	}
}

// extractKey извлекает ключ для сортировки по колонке k
func extractKey(line string, k int) string {
	if k <= 0 {
		return line
	}
	fields := strings.Split(line, "\t")
	if k <= len(fields) {
		return fields[k-1]
	}
	return ""
}

// isSorted проверяет, отсортированы ли строки
func isSorted(lines []string, cmp func(string, string) bool) bool {
	for i := 1; i < len(lines); i++ {
		if cmp(lines[i], lines[i-1]) {
			return false
		}
	}
	return true
}
