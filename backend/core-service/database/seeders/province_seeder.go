package seeders

import (
	"core-service/models"
	"core-service/pkg/utils"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

func RunProvinceSeeder(db *gorm.DB) error {
	basePath := getBasePath()

	filePath := filepath.Join(basePath, "database", "refs", "provinces.csv")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("province seeder: warning - file not found at %s, skipping\n", filePath)
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("province seeder: failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("province seeder: failed to read header: %w", err)
	}

	rowNum := 0
	for {
		record, err := reader.Read()
		if err != nil {
			// EOF is normal, break out
			break
		}

		rowNum++
		if len(record) < 2 {
			continue // Skip malformed rows
		}

		id, err := utils.StringToUint(record[0])
		if err != nil {
			// Log and continue with next row instead of breaking
			continue
		}

		province := models.Province{
			ID:   id,
			Name: record[1],
		}

		// Use FirstOrCreate to handle duplicates
		if err := db.Where("id = ?", id).FirstOrCreate(&province).Error; err != nil {
			return fmt.Errorf("province seeder: failed to insert row %d: %w", rowNum, err)
		}
	}

	return nil
}