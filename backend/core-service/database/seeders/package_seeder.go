package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PackageSeeder seeds sample package data with benefits
func RunPackageSeeder(db *gorm.DB) error {
	// Create Bronze Package
	bronzeID := uuid.New()
	bronze := models.Package{
		ID:              bronzeID,
		Name:            "Bronze Package",
		Description:     "Perfect for beginners looking to get started with driving basics.",
		PackageType:     models.PackageTypeBronze,
		Price:           1500000,
		DiscountPrice:   1250000,
		DurationMinutes: 60,
		TotalSessions:   6,
		Status:          models.PackageStatusActive,
		ImageURL:        "https://example.com/packages/bronze.jpg",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	bronzeBenefits := []models.PackageBenefit{
		{
			ID:        uuid.New(),
			PackageID: bronzeID,
			Title:     "4 Sessions with Instructor",
			Description: "Learn basics with experienced instructors",
			Icon:      "book-open",
			SortOrder: 1,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: bronzeID,
			Title:     "Basic Traffic Rules",
			Description: "Comprehensive traffic rules and signs training",
			Icon:      "traffic-cone",
			SortOrder: 2,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: bronzeID,
			Title:     "90-Minute Per Session",
			Description: "Adequate time for practical learning",
			Icon:      "clock",
			SortOrder: 3,
			CreatedAt: time.Now(),
		},
	}

	// Create Silver Package
	silverID := uuid.New()
	silver := models.Package{
		ID:              silverID,
		Name:            "Silver Package",
		Description:     "Our most popular package for comprehensive learning",
		PackageType:     models.PackageTypeSilver,
		Price:           1950000,
		DiscountPrice:   1600000,
		DurationMinutes: 60,
		TotalSessions:   8,
		Status:          models.PackageStatusActive,
		ImageURL:        "https://example.com/packages/silver.jpg",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	silverBenefits := []models.PackageBenefit{
		{
			ID:        uuid.New(),
			PackageID: silverID,
			Title:     "8 Sessions with Instructor",
			Description: "More practice time with professional instructors",
			Icon:      "users",
			SortOrder: 1,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: silverID,
			Title:     "City Driving Skills",
			Description: "Master driving in urban environments",
			Icon:      "building",
			SortOrder: 2,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: silverID,
			Title:     "Parking Practice",
			Description: "Learn parallel parking and perpendicular parking",
			Icon:      "car",
			SortOrder: 3,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: silverID,
			Title:     "Night Driving Training",
			Description: "Safe night driving techniques",
			Icon:      "moon",
			SortOrder: 4,
			CreatedAt: time.Now(),
		},
	}

	// Create Gold Package
	goldID := uuid.New()
	gold := models.Package{
		ID:              goldID,
		Name:            "Gold Package",
		Description:     "Complete mastery with unlimited support and premium perks",
		PackageType:     models.PackageTypeGold,
		Price:           2250000,
		DiscountPrice:   1850000,
		DurationMinutes: 60,
		TotalSessions:   10,
		Status:          models.PackageStatusActive,
		ImageURL:        "https://example.com/packages/gold.jpg",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	goldBenefits := []models.PackageBenefit{
		{
			ID:        uuid.New(),
			PackageID: goldID,
			Title:     "10 Sessions with Instructor",
			Description: "Extensive practice with certified instructors",
			Icon:      "award",
			SortOrder: 1,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: goldID,
			Title:     "Highway Driving",
			Description: "Learn safe highway driving techniques",
			Icon:      "road",
			SortOrder: 2,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: goldID,
			Title:     "Defensive Driving",
			Description: "Master defensive driving techniques",
			Icon:      "shield",
			SortOrder: 3,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: goldID,
			Title:     "Exam Preparation",
			Description: "Complete preparation for driving test",
			Icon:      "clipboard-check",
			SortOrder: 4,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: goldID,
			Title:     "Free Study Materials",
			Description: "Access to digital learning resources",
			Icon:      "book",
			SortOrder: 5,
			CreatedAt: time.Now(),
		},
	}

	// Create Platinum Package
	platinumID := uuid.New()
	platinum := models.Package{
		ID:              platinumID,
		Name:            "Platinum Package",
		Description:     "Complete mastery with unlimited support and premium perks",
		PackageType:     models.PackageTypePlatinum,
		Price:           2650000,
		DiscountPrice:   2100000,
		DurationMinutes: 60,
		TotalSessions:   12,
		Status:          models.PackageStatusActive,
		ImageURL:        "https://example.com/packages/platinum.jpg",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	platinumBenefits := []models.PackageBenefit{
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "12 Sessions 1-on-1 Coaching",
			Description: "Personalized instruction with senior instructors",
			Icon:      "user-check",
			SortOrder: 1,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "All Driving Scenarios",
			Description: "City, highway, mountain, and night driving",
			Icon:      "map",
			SortOrder: 2,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "Guaranteed Pass",
			Description: "If you don't pass, get free retake",
			Icon:      "check-circle",
			SortOrder: 3,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "Priority Scheduling",
			Description: "Flexible booking with priority slots",
			Icon:      "calendar",
			SortOrder: 4,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "Exclusive Study Materials",
			Description: "Premium video lessons and mock tests",
			Icon:      "video",
			SortOrder: 5,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			PackageID: platinumID,
			Title:     "Post-Training Support",
			Description: "3 months mentorship after completion",
			Icon:      "heart",
			SortOrder: 6,
			CreatedAt: time.Now(),
		},
	}

	// Create all packages with benefits
	packages := []struct {
		pkg     models.Package
		benefits []models.PackageBenefit
	}{
		{pkg: bronze, benefits: bronzeBenefits},
		{pkg: silver, benefits: silverBenefits},
		{pkg: gold, benefits: goldBenefits},
		{pkg: platinum, benefits: platinumBenefits},
	}

	for _, p := range packages {
		// Upsert package
		result := db.Where("id = ?", p.pkg.ID).FirstOrCreate(&p.pkg)
		if result.Error != nil {
			return result.Error
		}

		// Create benefits
		for _, benefit := range p.benefits {
			benefitResult := db.Where("package_id = ? AND title = ?", p.pkg.ID, benefit.Title).FirstOrCreate(&benefit)
			if benefitResult.Error != nil {
				return benefitResult.Error
			}
		}
	}

	return nil
}