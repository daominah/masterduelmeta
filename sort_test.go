package masterduelmeta

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func TestGetNextMonth(t *testing.T) {
	tests := []struct {
		month    string
		expected string
	}{
		{
			month:    "2024-01",
			expected: "2024-02",
		},
		{
			month:    "2024-12",
			expected: "2025-01",
		},
	}

	for _, tt := range tests {
		if got, err := GetNextMonth(tt.month); err != nil || got != tt.expected {
			t.Errorf(`GetNextMonth("%v") = %v, %v, want %v, nil`, tt.month, got, err, tt.expected)
		}
	}
}

func TestMonthsDecksToCSV(t *testing.T) {
	monthsDecks := map[string][]KeyCount[Archetype]{
		"2024-09": {
			{Key: "Voiceless Voice", Count: 40},
			{Key: "Yubel", Count: 30},
			{Key: "Tearlaments", Count: 20},
			{Key: "Centur-Ion", Count: 10},
		},
		"2024-10": {
			{Key: "Tenpai Dragon", Count: 50},
			{Key: "Yubel", Count: 20},
			{Key: "Voiceless Voice", Count: 20},
			{Key: "Mathmech", Count: 5},
			{Key: "Ritual Beasts", Count: 5},
		},
	}
	rows := MarshalMonthsDecksToCSV(monthsDecks)
	if len(rows) != 7 {
		t.Fatalf("len(rows) = %v, want 7", len(rows))
	}
	if got, want := strings.Join(rows[2], ","), `Voiceless Voice,40,40.0%,Tenpai Dragon,50,50.0%`; got != want {
		t.Errorf(`got rows[2] = %v, want %v`, got, want)
	}
	if got, want := strings.Join(rows[6], ","), `,,,Ritual Beasts,5,5.0%`; got != want {
		t.Errorf(`got rows[6] = %v, want %v`, got, want)
	}

	if true { // WriteAll csvData to a file
		file, err := os.Create("test.csv")
		if err != nil {
			t.Fatalf("error os.Create: %v", err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		writer.WriteAll(rows)
	}
}
