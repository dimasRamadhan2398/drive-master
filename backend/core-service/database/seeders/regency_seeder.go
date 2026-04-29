package seeders

import (
	"core-service/models"
	"core-service/pkg/utils"
	"encoding/csv"
	"os"

	"gorm.io/gorm"
)

func RunRegencySeeder(db *gorm.DB) error {
	file, err := os.Open("../refs/regencies.csv")
	if err != nil {
		panic(err)
	}
	
	defer file.Close()
	
	reader := csv.NewReader(file)
	
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		id, err := utils.StringToUint(record[0])
		if err != nil {
			break;
		}

		provinceId, err := utils.StringToUint(record[1])
		if err != nil {
			provinceId = 11
		}

		regency := models.Regency{
			ID: id,
			ProvinceID: uint(provinceId),
			Name: record[2],
			Type: record[3],
		}

		if err := db.Create(&regency).Error; err != nil {
			return err;
		}
	}

	return nil;
}