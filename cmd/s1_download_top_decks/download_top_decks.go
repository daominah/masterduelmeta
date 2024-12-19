package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	outputDir := filepath.Join(projectRoot, "downloaded_data")
	if _, err := os.Stat(outputDir); err != nil {
		log.Fatalf("error os.Stat outputDir: %v", err)
	}

	for _, month := range []string{
		"2022-01", "2022-02", "2022-03", "2022-04", "2022-05", "2022-06", "2022-07", "2022-08", "2022-09", "2022-10", "2022-11", "2022-12",
		"2023-01", "2023-02", "2023-03", "2023-04", "2023-05", "2023-06", "2023-07", "2023-08", "2023-09", "2023-10", "2023-11", "2023-12",
		"2024-01", "2024-02", "2024-03", "2024-04", "2024-05", "2024-06", "2024-07", "2024-08", "2024-09", "2024-10",
		"2024-11", "2024-12",
	} {
		log.Printf("downloading data for %v", month)
		outputPath := filepath.Join(outputDir, fmt.Sprintf(`decks_%v.json`, month))
		_, err := os.Stat(outputPath)
		if err == nil {
			log.Printf("file %v already exists", outputPath)
			continue
		}

		file, err := os.Create(outputPath)
		if err != nil {
			log.Printf("error os.Create: %v", err)
			continue
		}
		data, err := DownloadTopDecks(month)
		if err != nil {
			log.Printf("error DownloadTopDecks: %v", err)
			continue
		}
		_, err = file.Write(data)
		if err != nil {
			log.Printf("error file.Write: %v", err)
			continue
		}
		log.Printf("downloaded %v, size %v KiB", outputPath, len(data)/1024)
	}
}

// DownloadTopDecks downloads top decks of a month from masterduelmeta.com,
// Arguments: - month: format "yyyy-mm", e.g. "2006-01"
func DownloadTopDecks(month string) ([]byte, error) {
	nextMonth, err := masterduelmeta.GetNextMonth(month)
	if err != nil {
		return nil, fmt.Errorf("GetNextMonth: %v", err)
	}
	u := fmt.Sprintf(`https://www.masterduelmeta.com/api/v1/top-decks`+
		`?created[$gte]=%v&created[$lt]=%v&fields=-_id,-__v,-notes&limit=0`, month, nextMonth)
	log.Printf("doing http.Get %v", u)
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad resp.StatusCode: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v", err)
	}

	return body, nil
}
