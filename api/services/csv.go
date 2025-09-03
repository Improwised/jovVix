package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateCsv[T any](filePath, id string, analysis T, getData func(T) [][]string) (string, error) {
	fullPath := filepath.Join(filePath, id+".csv")

	csvFile, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	data := getData(analysis)
	if data == nil {
		return "", fmt.Errorf("no data found")
	}

	for _, row := range data {
		if err := csvwriter.Write(row); err != nil {
			return "", err
		}
	}

	if err = csvwriter.Error(); err != nil {
		return "", err
	}

	return fullPath, nil
}
