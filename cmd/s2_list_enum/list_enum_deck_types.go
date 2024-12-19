package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/daominah/masterduelmeta"
	"github.com/mywrap/gofast"
)

func main() {
	projectRoot, err := gofast.GetProjectRootPath()
	if err != nil {
		log.Fatalf("error gofast.GetProjectRootPath: %v", err)
	}
	dataDir := filepath.Join(projectRoot, "downloaded_data")
	log.Printf("dataDir: %v", dataDir)
	if _, err := os.Stat(dataDir); err != nil {
		log.Fatalf("error os.Stat dataDir: %v", err)
	}

	var allDecks []masterduelmeta.Deck
	for year := 2022; year <= 2025; year++ {
		for month := 1; month <= 12; month++ {
			monthStr := fmt.Sprintf("%v-%02v", year, month)
			dataPath := filepath.Join(dataDir, fmt.Sprintf("decks_%v.json", monthStr))
			data, err := os.ReadFile(dataPath)
			if err != nil {
				log.Printf("error os.ReadFile: %v", err)
				if year >= 2025 {
					break
				}
				continue
			}
			decks, err := masterduelmeta.ParseDecks(data)
			if err != nil {
				log.Printf("error masterduelmeta.ParseDecks: %v", err)
				continue
			}
			//log.Printf("month %v: len(decks): %v", monthStr, len(decks))
			allDecks = append(allDecks, decks...)
		}
	}

	countDeckTypes := make(map[string]int)
	countRankTypes := make(map[masterduelmeta.RankedType]int)
	countTournamentTypes := make(map[masterduelmeta.TournamentType]int)
	countEngines := make(map[string]int)
	for _, deck := range allDecks {
		countDeckTypes[deck.DeckType.Name] += 1
		countRankTypes[deck.RankedType.Name]++
		countTournamentTypes[deck.TournamentType.Name]++
		for _, engine := range deck.Engines {
			countEngines[engine.Name]++
		}
	}

	log.Printf("________________________________________________________")
	sortedRankTypes := masterduelmeta.SortMapByValueDesc(countRankTypes)
	for i, rankType := range sortedRankTypes {
		log.Printf("rankType %03v: %40v: %v", i, rankType.Key, rankType.Count)
	}

	log.Printf("________________________________________________________")
	sortedTournamentTypes := masterduelmeta.SortMapByValueDesc(countTournamentTypes)
	for i, tournamentType := range sortedTournamentTypes {
		log.Printf("tournamentType %03v: %40v: %v", i, tournamentType.Key, tournamentType.Count)
	}

	//log.Printf("________________________________________________________")
	//sortedDeckTypes := masterduelmeta.SortMapByValueDesc(countDeckTypes)
	//for i, deckType := range sortedDeckTypes {
	//	log.Printf("deckType %03v: %40v: %v", i, deckType.Key, deckType.Count)
	//}

	//log.Printf("________________________________________________________")
	//sortedEngines := masterduelmeta.SortMapByValueDesc(countEngines)
	//for i, engine := range sortedEngines {
	//	log.Printf("engine %03v: %40v: %v", i, engine.Key, engine.Count)
	//}
}
