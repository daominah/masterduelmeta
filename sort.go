package masterduelmeta

import (
	"fmt"
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

func MarshalMonthsDecksToCSV(monthsDecks map[string][]KeyCount[Archetype]) [][]string {
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
			percent := fmt.Sprintf("%.1f", float64(deckCount.Count)/float64(sumDecksCount)*100)
			column2 = append(column2, percent+`%`)
		}
		columns = append(columns, column0, column1, column2)
	}
	return RotateMatrix(columns)
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
