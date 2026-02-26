package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Row хранит исходную строку и подготовленные данные для быстрого сравнения
type Row struct {
	Content string
	Column  string
	Number  float64
	Month   int
}

// Config содержит настройки из флагов командной строки
type Config struct {
	k int  // колонка
	n bool // числа
	r bool // реверс
	u bool // уникальные
	M bool // месяцы
	b bool // игнорировать пробелы в конце
	c bool // только проверка сортировки
	h bool // человекочитаемые числа
}

func main() {
	cfg := parseFlags()

	rows, err := readLines(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Если включен флаг -c (проверка отсортированности)
	if cfg.c {
		if isSorted(rows, cfg) {
			fmt.Println("Data is already sorted")
		} else {
			fmt.Println("Data is NOT sorted")
		}
		return
	}

	// Основная сортировка
	sortRows(rows, cfg)

	// Обработка уникальности (-u)
	if cfg.u {
		rows = unique(rows)
	}

	// Вывод результата
	for _, row := range rows {
		fmt.Println(row.Content)
	}
}

// parseFlags инициализирует и разбирает аргументы командной строки
func parseFlags() Config {
	cfg := Config{}
	flag.IntVar(&cfg.k, "k", 1, "sort by column number")
	flag.BoolVar(&cfg.n, "n", false, "sort by numeric value")
	flag.BoolVar(&cfg.r, "r", false, "reverse order")
	flag.BoolVar(&cfg.u, "u", false, "unique lines only")
	flag.BoolVar(&cfg.M, "M", false, "sort by month")
	flag.BoolVar(&cfg.b, "b", false, "ignore trailing blanks")
	flag.BoolVar(&cfg.c, "c", false, "check if sorted")
	flag.BoolVar(&cfg.h, "h", false, "sort by human readable numbers")
	flag.Parse()
	return cfg
}

// readLines читает данные и сразу подготавливает ключи сортировки в Row
func readLines(cfg Config) ([]Row, error) {
	var rows []Row
	var input *os.File = os.Stdin

	if len(flag.Args()) > 0 {
		f, err := os.Open(flag.Args()[0])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		input = f
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.b {
			line = strings.TrimRight(line, " ")
		}

		col := getColumn(line, cfg.k)
		row := Row{
			Content: line,
			Column:  col,
		}

		// Пре-парсинг для ускорения сортировки
		if cfg.n || cfg.h {
			row.Number = parseNumeric(col, cfg.h)
		}
		if cfg.M {
			row.Month = parseMonth(col)
		}

		rows = append(rows, row)
	}
	return rows, scanner.Err()
}

// sortRows — сердце программы, реализует логику сравнения
func sortRows(rows []Row, cfg Config) {
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		var less bool

		switch {
		case cfg.M:
			less = a.Month < b.Month
		case cfg.n || cfg.h:
			less = a.Number < b.Number
		default:
			less = a.Column < b.Column
		}

		if cfg.r {
			return !less
		}
		return less
	})
}

// getColumn извлекает N-ую колонку (разделитель таб)
func getColumn(s string, k int) string {
	parts := strings.Split(s, "\t")
	if k > len(parts) || k <= 0 {
		return ""
	}
	return parts[k-1]
}

// parseNumeric обрабатывает числа и суффиксы (K, M, G)
func parseNumeric(s string, human bool) float64 {
	if s == "" {
		return 0
	}

	multiplier := 1.0
	if human {
		s = strings.ToUpper(s)
		switch {
		case strings.HasSuffix(s, "K"):
			multiplier = 1024
			s = s[:len(s)-1]
		case strings.HasSuffix(s, "M"):
			multiplier = 1024 * 1024
			s = s[:len(s)-1]
		case strings.HasSuffix(s, "G"):
			multiplier = 1024 * 1024 * 1024
			s = s[:len(s)-1]
		}
	}

	val, _ := strconv.ParseFloat(s, 64)
	return val * multiplier
}

// parseMonth превращает название месяца в число 1-12
func parseMonth(s string) int {
	months := map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6,
		"JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
	}
	val, ok := months[strings.ToUpper(s[:3])]
	if !ok {
		return 0
	}
	return val
}

// unique удаляет дубликаты из отсортированного слайса
func unique(rows []Row) []Row {
	if len(rows) <= 1 {
		return rows
	}
	result := []Row{rows[0]}
	for i := 1; i < len(rows); i++ {
		if rows[i].Content != rows[i-1].Content {
			result = append(result, rows[i])
		}
	}
	return result
}

func isSorted(rows []Row, cfg Config) bool {
	return sort.SliceIsSorted(rows, func(i, j int) bool {
		return true
	})
}
