package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mywrap/gofast"
)

func TestDownloadTopDecks(t *testing.T) {
	month := "2024-12"
	t.Skip("this test always redownload data and overwrite the existing data for " + month)
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
}
