package main

import (
	"encoding/csv"
	"log/slog"
	"os"
)

func CreateStore() {
	file, err := os.Create(storeFile)
	if err != nil {
		slog.Error("Failed to create store file")
		return
	}

	_ = file.Close()
}

func LoadStore() map[string]string {
	file, err := os.Open(storeFile)
	if err != nil {
		slog.Info("No store file found; creating new")
		CreateStore()
		file, _ = os.Open(storeFile)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		slog.Error("Failed to parse CSV")
		return nil
	}

	return RecordsToMap(records)
}

func RecordsToMap(records [][]string) map[string]string {
	slugs := map[string]string{}

	for rowIndex, row := range records {
		// should be even # of items in the row to interpret as key value
		if len(row)%2 != 0 {
			slog.Error("Row has wrong format", "row", rowIndex, "row data", row)
			return nil
		}

		cursor := 0
		key := ""

		for _, cell := range row {
			switch cursor {
			case 0:
				key = cell
				cursor++
			case 1:
				// we have key and value now, insert it into map
				slugs[key] = cell
				cursor = 0
			}
		}

	}

	return slugs
}

func MapToRecords(slugs map[string]string) [][]string {
	var records []string

	for slug := range slugs {
		key := slug
		value := slugs[key]

		records = append(records, key)
		records = append(records, value)
	}

	return [][]string{records}
}

func SaveStore(slugs map[string]string) {
	records := MapToRecords(slugs)

	file, err := os.Create(storeFile)
	if err != nil {
		slog.Error("Failed to create store file")
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		slog.Error("Failed to save records to store")
		return
	}
}
