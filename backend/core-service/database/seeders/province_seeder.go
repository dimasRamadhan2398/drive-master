package seeders

import (
	"core-service/models"
	"core-service/pkg/utils"
	"encoding/csv"
	"os"

	"gorm.io/gorm"
)

func RunProvinceSeeder(db *gorm.DB) error {
	file, err := os.Open("../refs/provinces.csv")
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

		id, err :=utils.StringToUint(record[0])
		if err != nil {
			break;
		}
		
		province := models.Province{
			ID: id,
			Name: record[1],
		}
		
		if err := db.Create(&province).Error; err != nil {
			return err;
		}
	}
	return nil
}