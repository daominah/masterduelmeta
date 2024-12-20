package masterduelmeta

import (
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

	t.Run("TestMarshalMonthsDecksToCSVGroupByMonth", func(t *testing.T) {
		rows := MarshalMonthsDecksToCSVGroupByMonth(monthsDecks)
		if len(rows) != 7 {
			t.Fatalf("len(rows) = %v, want 7", len(rows))
		}
		if got, want := strings.Join(rows[2], ","), `Voiceless Voice,40,40.0%,Tenpai Dragon,50,50.0%`; got != want {
			t.Errorf(`got rows[2] = %v, want %v`, got, want)
		}
		if got, want := strings.Join(rows[6], ","), `,,,Ritual Beasts,5,5.0%`; got != want {
			t.Errorf(`got rows[6] = %v, want %v`, got, want)
		}

		if false {
			WriteTestCSVFile(rows) //
		}
	})

	t.Run("TestMarshalMonthsDecksToCSV", func(t *testing.T) {
		rows := MarshalMonthsDecksToCSV(monthsDecks)
		if len(rows) != 10 {
			t.Fatalf("len(rows) = %v, want 3", len(rows))
		}
		if got, want := strings.Join(rows[0], ","), `Month,Deck,Percent`; got != want {
			t.Errorf(`got rows[0] = %v, want %v`, got, want)
		}
		if got, want := strings.Join(rows[1], ","), `2024-09,Voiceless Voice,40.0`; got != want {
			t.Errorf(`got rows[1] = %v, want %v`, got, want)
		}
		if got, want := strings.Join(rows[5], ","), `2024-10,Tenpai Dragon,50.0`; got != want {
			t.Errorf(`got rows[2] = %v, want %v`, got, want)
		}
		if got, want := strings.Join(rows[9], ","), `2024-10,Ritual Beasts,5.0`; got != want {
			t.Errorf(`got rows[3] = %v, want %v`, got, want)
		}

		if false {
			WriteTestCSVFile(rows) //
		}
	})
}
