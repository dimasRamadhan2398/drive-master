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

func RunDistrictsSeeder(db *gorm.DB) error {
	basePath := getBasePath()

	filePath := filepath.Join(basePath, "database", "refs", "districts.csv")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("district seeder: file not found at %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("district seeder: failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("district seeder: failed to read header: %w", err)
	}

	rowNum := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		rowNum++
		if len(record) < 2 {
			continue // Skip malformed rows
		}

		regencyId, err := utils.StringToUint(record[1])
		if err != nil {
			continue
		}

		district := models.District{
			Name:      record[0],
			RegencyID: regencyId,
		}

		// Use FirstOrCreate to handle duplicates
		if err := db.Where("name = ? AND regency_id = ?", record[0], regencyId).FirstOrCreate(&district).Error; err != nil {
			return fmt.Errorf("district seeder: failed to insert row %d: %w", rowNum, err)
		}
	}

	return nil
}