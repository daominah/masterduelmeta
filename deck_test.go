package masterduelmeta

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/mywrap/gofast"
)

//go:embed deck_test.json
var testParseDecksData []byte

func TestParseDecks(t *testing.T) {
	// DaoMinAh decks https://www.masterduelmeta.com/user/646456f28d4d91b234ae3a57

	if len(testParseDecksData) == 0 {
		t.Fatal("test2DecksData is empty, check embed file path")
	}
	decks, err := ParseDecks(testParseDecksData)
	if err != nil {
		t.Fatalf("error ParseDecks: %v", err)
	}
	if len(decks) != 11 {
		t.Fatalf("len(decks) got %v, want 2", len(decks))
	}

	countBlueEyesDecks := 0
	for _, deck := range decks {
		if deck.DeckType.Name == "Blue-Eyes" {
			countBlueEyesDecks++
		}
	}
	if countBlueEyesDecks != 3 {
		t.Fatalf("countBlueEyesDecks got %v, want 3", countBlueEyesDecks)
	}

	//for _, deck := range decks {
	//	t.Logf("deck: %v", deck.DeckType)
	//}
}

func TestAnalyze_2024_06(t *testing.T) {
	projectRoot, _ := gofast.GetProjectRootPath()
	deckData, err := os.ReadFile(filepath.Join(projectRoot, "downloaded_data", "decks_2024-06.json"))
	if err != nil || len(deckData) == 0 {
		t.Fatalf("error ParseDecks: %v", err)
	}
	decks, err := ParseDecks(deckData)
	if err != nil {
		t.Fatalf("error ParseDecks: %v", err)
	}
	t.Logf("len(decks): %v", len(decks)) // Output: len(decks): 1511

	countRankTypes := make(map[any]int)
	for _, deck := range decks {
		countRankTypes[deck.RankedType.Name]++
	}
	sortedRankTypes := SortMapByValueDesc(countRankTypes)
	t.Logf("rank types:")
	for _, rankTypeCount := range sortedRankTypes {
		t.Logf("%20v: %v", rankTypeCount.Key, rankTypeCount.Count)
	}

	// deck type can be "Tenpai Dragon", "Voiceless Voice", "Yubel", ...
	countDeckTypes := make(map[any]int)
	for _, deck := range decks {
		countDeckTypes[deck.DeckType.Name]++
	}
	sortedDeckTypes := SortMapByValueDesc(countDeckTypes)
	t.Logf("___________________________________________________________")
	t.Logf("deck types:")
	for _, deckTypeCount := range sortedDeckTypes {
		t.Logf("%20v: %v", deckTypeCount.Key, deckTypeCount.Count)
	}

	for _, deck := range decks {
		if deck.RankedType.Name == "" {
			t.Logf("deck: %#v", deck)
			break
		}
	}
}
