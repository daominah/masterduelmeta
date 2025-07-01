package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

	beginMonth := "2024-12"
	endMonth := "2025-06"

	startTime, err := time.Parse("2006-01", beginMonth)
	if err != nil {
		log.Fatalf("error parsing beginMonth: %v", err)
	}
	endTime, err := time.Parse("2006-01", endMonth)
	if err != nil {
		log.Fatalf("error parsing endMonth: %v", err)
	}

	monthsRankDecks := make(map[string][]ygo.KeyCount[ygo.Archetype])
	var allRankDecks []ygo.Deck

	for current := startTime; !current.After(endTime); current = current.AddDate(0, 1, 0) {
		monthStr := current.Format("2006-01")
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
				allRankDecks = append(allRankDecks, deck)
				archetypesCount[deck.Archetype()]++
			}
		}
		monthsRankDecks[monthStr] = ygo.SortMapByValueDesc(archetypesCount)
	}

	log.Printf("________________________________")

	monthDecksColumnFile := fmt.Sprintf("cmd/s3_analyze_to_visualize/ranked_decks_%s_to_%s.csv", beginMonth, endMonth)
	timeSeriesFile := fmt.Sprintf("visualization/time_series_%s_to_%s.csv", beginMonth, endMonth)

	for outputFilePath, marshalFunc := range map[string]func(map[string][]ygo.KeyCount[ygo.Archetype]) [][]string{
		monthDecksColumnFile: ygo.MarshalMonthsDecksToCSVGroupByMonth,
		timeSeriesFile:       ygo.MarshalMonthsDecksToCSV,
	} {
		outputFilePath = filepath.Join(projectRoot, "", outputFilePath)
		file, err := os.Create(outputFilePath)
		if err != nil {
			log.Fatalf("error os.Create: %v", err)
		}
		csvWriter := csv.NewWriter(file)
		csvWriter.WriteAll(marshalFunc(monthsRankDecks))
		csvWriter.Flush()
		file.Close()
		log.Printf("wrote %v", outputFilePath)
	}

	if false { // specific query to clarify deck type
		log.Printf("________________________________")
		for _, deck := range allRankDecks {
			if deck.CheckContainsCard("Qebehsenuef, Protection of Horus") &&
				// deck.CheckContainsCard("Secret Village of the Spellcasters") &&
				deck.CheckContainsCard("Summon Limit") {
				log.Printf("deck Qebehsenuef type %v", deck.DeckType.Name)
			}
		}
	}
}
