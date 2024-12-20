package masterduelmeta

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type KeyCount[K comparable] struct {
	Key   K
	Count int
}

type SortByValue[K comparable] []KeyCount[K]

func (a SortByValue[K]) Len() int           { return len(a) }
func (a SortByValue[K]) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByValue[K]) Less(i, j int) bool { return a[i].Count > a[j].Count }

func SortMapByValueDesc[K comparable](mapDeckTypes map[K]int) []KeyCount[K] {
	var deckTypesCount []KeyCount[K]
	for deckType, count := range mapDeckTypes {
		deckTypesCount = append(deckTypesCount, KeyCount[K]{deckType, count})
	}
	sort.Sort(SortByValue[K](deckTypesCount))
	return deckTypesCount
}

// GetNextMonth receives a month string in the format "yyyy-mm",
// returns the next month in the same format,
// example: GetNextMonth("2024-01") = "2024-02", GetNextMonth("2024-12") = "2025-01".
func GetNextMonth(month string) (string, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return "", fmt.Errorf(`time.Parse: %v`, err)
	}
	nextMonth := t.AddDate(0, 1, 0)
	return nextMonth.Format("2006-01"), nil
}

func MarshalMonthsDecksToCSVGroupByMonth(monthsDecks map[string][]KeyCount[Archetype]) [][]string {
	var months []string
	for month := range monthsDecks {
		months = append(months, month)
	}
	sort.Strings(months)

	var columns [][]string
	for _, month := range months {
		sumDecksCount := 0
		for _, deckCount := range monthsDecks[month] {
			sumDecksCount += deckCount.Count
		}
		column0 := []string{month, strconv.Itoa(sumDecksCount)} // deck name
		column1 := []string{"", ""}                             // count
		column2 := []string{"", ""}                             // percentage
		for _, deckCount := range monthsDecks[month] {
			column0 = append(column0, string(deckCount.Key))
			column1 = append(column1, strconv.Itoa(deckCount.Count))
			percent := float64(deckCount.Count) / float64(sumDecksCount) * 100
			column2 = append(column2, fmt.Sprintf("%.1f", percent)+`%`)
		}
		columns = append(columns, column0, column1, column2)
	}
	return RotateMatrix(columns)
}

func MarshalMonthsDecksToCSV(monthsDecks map[string][]KeyCount[Archetype]) [][]string {
	var months []string
	for month := range monthsDecks {
		months = append(months, month)
	}
	sort.Strings(months)

	rows := [][]string{{"Month", "Deck", "Percent"}}
	for _, month := range months {
		sumDecksCount := 0
		for _, deckCount := range monthsDecks[month] {
			sumDecksCount += deckCount.Count
		}
		for _, deckCount := range monthsDecks[month] {
			percent := float64(deckCount.Count) / float64(sumDecksCount) * 100
			rows = append(rows, []string{
				month,
				string(deckCount.Key),
				fmt.Sprintf("%.1f", percent),
			})
		}
	}
	return rows
}

// RotateMatrix rotate rows and columns of a matrix,
// can handle a matrix with different row lengths without panicking.
func RotateMatrix(matrix [][]string) [][]string {
	if len(matrix) == 0 {
		return [][]string{}
	}
	// determine the maximum length of the rows
	maxLen := 0
	for _, row := range matrix {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	rotated := make([][]string, maxLen)
	for i := range rotated {
		rotated[i] = make([]string, len(matrix))
	}
	for i, row := range matrix {
		for j, val := range row {
			rotated[j][i] = val
		}
	}

	return rotated
}

// WriteTestCSVFile writes data to "test_output.csv" (at the current running directory),
func WriteTestCSVFile(rows [][]string) error {
	file, err := os.Create("test_output.csv")
	if err != nil {
		return fmt.Errorf(`os.Create: %v`, err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.WriteAll(rows)
	if err != nil {
		return fmt.Errorf(`writer.WriteAll: %v`, err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf(`writer.Error: %v`, err)
	}
	return nil
}
