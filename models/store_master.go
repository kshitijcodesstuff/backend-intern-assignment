package models

import (
	"encoding/csv"
	"os"
)

var storeMaster = map[string]bool{}

// LoadStoreMaster preloads StoreMaster.csv
func LoadStoreMaster(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, record := range records[1:] { // Skip header row
		storeID := record[2]
		storeMaster[storeID] = true
	}
}

// IsValidStore checks if a store ID exists in the master list
func IsValidStore(storeID string) bool {
	return storeMaster[storeID]
}

func InitTestStoreMaster() {
	storeMaster = map[string]bool{
		"RP00001": true,
		"RP00002": true,
	}
}
