package seeders

import (
	"core-service/models"
	"core-service/pkg/utils"
	"encoding/csv"
	"os"

	"gorm.io/gorm"
)

func RunDistrictsSeeder(db *gorm.DB) error {
	file, err := os.Open("../refs/districts.csv")
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

		regencyId, err := utils.StringToUint(record[1])
		
		district := models.District{
			Name: record[0],
			RegencyID: regencyId,
		}
		
		if err := db.Create(&district).Error; err != nil {
			return err;
		}
	}
	return nil
}