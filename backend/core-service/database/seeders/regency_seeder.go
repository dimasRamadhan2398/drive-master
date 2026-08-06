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

func RunRegencySeeder(db *gorm.DB) error {
	basePath := getBasePath()

	filePath := filepath.Join(basePath, "database", "refs", "regencies.csv")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("regency seeder: warning - file not found at %s, skipping\n", filePath)
		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("regency seeder: failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("regency seeder: failed to read header: %w", err)
	}

	rowNum := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		rowNum++
		if len(record) < 4 {
			continue // Skip malformed rows
		}

		id, err := utils.StringToUint(record[0])
		if err != nil {
			continue
		}

		provinceId, err := utils.StringToUint(record[1])
		if err != nil {
			provinceId = 0
		}

		regency := models.Regency{
			ID:         id,
			ProvinceID: provinceId,
			Name:       record[2],
			Type:       record[3],
		}

		// Use FirstOrCreate to handle duplicates
		if err := db.Where("id = ?", id).FirstOrCreate(&regency).Error; err != nil {
			return fmt.Errorf("regency seeder: failed to insert row %d: %w", rowNum, err)
		}
	}

	return nil
}