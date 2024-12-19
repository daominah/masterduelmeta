package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	var allRankDecks []ygo.Deck
	monthsRankDecks := make(map[string][]ygo.KeyCount[ygo.Archetype])

	year := 2024
	for month := 1; month <= 12; month++ {
		monthStr := fmt.Sprintf("%v-%02v", year, month)
		dataPath := filepath.Join(dataDir, fmt.Sprintf("decks_%v.json", monthStr))
		data, err := os.ReadFile(dataPath)
		if err != nil {
			log.Printf("error os.ReadFile: %v", err)
			continue
		}
		decks, err := ygo.ParseDecks(data)
		if err != nil {
			log.Printf("error ygo.ParseDecks: %v", err)
			continue
		}
		var rankedDecks []ygo.Deck
		archetypesCount := make(map[ygo.Archetype]int)
		for _, deck := range decks {
			if deck.CheckIsNormalRank() {
				rankedDecks = append(rankedDecks, deck)
				archetypesCount[deck.Archetype()]++
			}
		}
		allRankDecks = append(allRankDecks, rankedDecks...)
		monthsRankDecks[monthStr] = ygo.SortMapByValueDesc(archetypesCount)
		if len(monthsRankDecks[monthStr]) > 10 {
			monthsRankDecks[monthStr] = monthsRankDecks[monthStr][:10]
		}
		log.Printf("month %v: len(decks) %v, len(rankedDecks) %v", monthStr, len(decks), len(rankedDecks))
	}

	outputFilePath := filepath.Join(projectRoot, "cmd/s3_analyze_2024/ranked_decks_2024_12.csv")
	file, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("error os.Create: %v", err)
	}
	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(ygo.MarshalMonthsDecksToCSV(monthsRankDecks))
	csvWriter.Flush()
	file.Close()
	log.Printf("wrote %v", outputFilePath)

	_ = allRankDecks
}
