package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	ygo "github.com/daominah/masterduelmeta"
	"github.com/mywrap/gofast"
)

func TestReDownloadTopDecksCurrentMonth(t *testing.T) {
	month := time.Now().Format("2006-01")
	// month, _ := "2024-12", time.Now

	// t.Skip("this test always redownload data and overwrite the existing data for " + month)

	t.Logf("redownload data for %v", month)
	data, err := DownloadTopDecks(month)
	if err != nil {
		t.Fatalf("error DownloadTopDecks: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("len(data) == 0")
	}
	projectRoot, _ := gofast.GetProjectRootPath()
	outputPath := filepath.Join(projectRoot, "downloaded_data", fmt.Sprintf(`decks_%v.json`, month))
	t.Logf("outputPath: %v", outputPath)
	err = os.WriteFile(outputPath, data, 0o666)
	if err != nil {
		t.Fatalf("error WriteFile: %v", err)
	}
	t.Logf("downloaded %v, size %v KiB", outputPath, len(data)/1024)

	decks, err := ygo.ParseDecks(data)
	if err != nil {
		t.Fatalf("error ygo.ParseDecks: %v", err)
	}
	t.Logf("len(decks) in %v: %v", len(decks), month)
}
