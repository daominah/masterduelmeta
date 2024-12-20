package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ygo "github.com/daominah/masterduelmeta"
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

	var allDecks []ygo.Deck
	for year := 2024; year <= 2025; year++ {
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
			decks, err := ygo.ParseDecks(data)
			if err != nil {
				log.Printf("error ygo.ParseDecks: %v", err)
				continue
			}
			// log.Printf("month %v: len(decks): %v", monthStr, len(decks))
			allDecks = append(allDecks, decks...)
		}
	}

	countDeckTypes := make(map[string]int)
	countDeckTypesInRank := make(map[string]int)
	countArchetypesInRank := make(map[ygo.Archetype]int)
	countRankTypes := make(map[ygo.RankedType]int)
	countTournamentTypes := make(map[ygo.TournamentType]int)
	countEngines := make(map[string]int)
	for _, deck := range allDecks {
		countDeckTypes[deck.DeckType.Name] += 1
		if deck.CheckIsNormalRank() {
			countDeckTypesInRank[deck.DeckType.Name] += 1
			countArchetypesInRank[deck.Archetype()] += 1
		}
		countRankTypes[deck.RankedType.Name]++
		countTournamentTypes[deck.TournamentType.Name]++
		for _, engine := range deck.Engines {
			countEngines[engine.Name]++
		}
	}

	log.Printf("________________________________________________________")
	sortedRankTypes := ygo.SortMapByValueDesc(countRankTypes)
	for i, rankType := range sortedRankTypes {
		log.Printf("rankType %03v: %40v: %v", i, rankType.Key, rankType.Count)
	}

	log.Printf("________________________________________________________")
	sortedTournamentTypes := ygo.SortMapByValueDesc(countTournamentTypes)
	for i, tournamentType := range sortedTournamentTypes {
		log.Printf("tournamentType %03v: %40v: %v", i, tournamentType.Key, tournamentType.Count)
	}

	var lines []string
	log.Printf("________________________________________________________")
	sortedDeckTypes := ygo.SortMapByValueDesc(countArchetypesInRank)
	for _, deckType := range sortedDeckTypes {
		line := fmt.Sprintf("%40v: %v", deckType.Key, deckType.Count)
		lines = append(lines, line)
		// fmt.Println(line)
	}
	archetypesOutFile := filepath.Join(projectRoot, "cmd/s2_list_enum", "archetypes.txt")
	err = os.WriteFile(archetypesOutFile, []byte(strings.Join(lines, "\n")), 0o666)
	if err != nil {
		log.Fatalf("error os.WriteFile: %v", err)
	}
	log.Printf("wrote %v", archetypesOutFile)

	//log.Printf("________________________________________________________")
	//sortedEngines := ygo.SortMapByValueDesc(countEngines)
	//for i, engine := range sortedEngines {
	//	log.Printf("engine %03v: %40v: %v", i, engine.Key, engine.Count)
	//}
}
