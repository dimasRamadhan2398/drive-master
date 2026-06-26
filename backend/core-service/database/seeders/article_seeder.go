package seeders

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// RunArticleSeeder seeds sample article data
func RunArticleSeeder(db *gorm.DB) error {
	articles := []models.Article{
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000001"),
			Title:          "5 Tips for First-Time Drivers to Pass Their SIM A Test",
			Slug:           "5-tips-first-time-drivers-pass-sim-a-test",
			LeadParagraph:  "Getting your SIM A license is a milestone. Here are expert tips to help you pass on your first attempt.",
			Content: `<p>Passing your SIM A driving test on the first try requires preparation and confidence. Here are five essential tips that have helped thousands of new drivers succeed.</p>

<h2>1. Practice Regularly</h2>
<p>Consistency is key. Aim to practice at least 3-4 times per week before your test date. Each session should be at least 60 minutes to build muscle memory.</p>

<h2>2. Master the Basics</h2>
<p>Before attempting complex maneuvers, ensure you have complete control of the basics: starting, steering, braking, and parking.</p>

<h2>3. Know Your Vehicle</h2>
<p>Understand your vehicle's dimensions, blind spots, and controls. This knowledge will help you navigate confidently during the test.</p>

<h2>4. Study Traffic Rules</h2>
<p>The written portion of the test requires knowledge of Indonesian traffic regulations. Study the rulebook thoroughly.</p>

<h2>5. Stay Calm and Focused</h2>
<p>Nervousness is normal but don't let it control you. Take deep breaths, follow instructions carefully, and trust your preparation.</p>`,
			Footer:           "Good luck with your test! Remember, practice makes perfect.",
			MetaTitle:         "5 Tips for First-Time Drivers to Pass Their SIM A Test",
			MetaDescription:   "Discover expert tips to help you pass your SIM A driving test on the first attempt. Essential advice for new drivers in Indonesia.",
			MetaKeywords:      "SIM A, driving test, first time driver, driving tips, Indonesia driver's license",
			CanonicalURL:      "https://drive-master.com/articles/5-tips-first-time-drivers-pass-sim-a-test",
			OGTitle:           "5 Tips to Pass Your SIM A Test First Time",
			OGDescription:     "Expert advice to help new drivers succeed in their SIM A test.",
			OGImage:           "https://example.com/images/sim-a-tips.jpg",
			FeaturedImage:     "https://example.com/images/driving-tips-hero.jpg",
			ThumbnailURL:      "https://example.com/images/driving-tips-thumb.jpg",
			Tags:              pq.StringArray{"tips", "sim-a", "beginners", "driving-test"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, -1, 0)),
			ViewCount:         1523,
			LikeCount:         89,
			ShareCount:        45,
			ReadingTime:       5,
			IsFeatured:        true,
			IsSpotlight:       false,
			Priority:          10,
			Language:          "en",
			Version:           1,
			EstimatedTime:     5,
			CTAText:           "Book Your First Lesson",
			CTALink:           "/packages",
			AvgReadTime:       4.5,
			BounceRate:        0.32,
			CompletionRate:    0.78,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000002"),
			Title:          "Understanding Indonesian Road Signs: A Complete Guide",
			Slug:           "understanding-indonesian-road-signs-complete-guide",
			LeadParagraph:  "Navigate Indonesian roads with confidence by mastering these essential traffic signs and their meanings.",
			Content: `<p>Indonesian roads can be challenging for new drivers. Understanding road signs is crucial for safe navigation.</p>

<h2>Warning Signs</h2>
<p>Warning signs in Indonesia are typically triangular with red borders. They alert drivers to potential hazards ahead.</p>

<h2>Prohibitory Signs</h2>
<p>These circular signs with red borders indicate actions that are not permitted, such as no entry or no left turn.</p>

<h2>Mandatory Signs</h2>
<p>Blue circular signs indicate mandatory actions, such as the direction you must take or activities you must do.</p>

<h2>Information Signs</h2>
<p>These rectangular blue or green signs provide helpful information about nearby facilities and directions.</p>`,
			Footer:           "Stay safe on Indonesian roads by always paying attention to traffic signs.",
			MetaTitle:         "Indonesian Road Signs Guide - Complete Reference",
			MetaDescription:   "Learn all Indonesian road signs with our comprehensive guide. Essential reading for new drivers.",
			MetaKeywords:      "Indonesian road signs, traffic signs, driving in Indonesia, road safety",
			Tags:              pq.StringArray{"road-safety", "traffic-rules", "beginners", "indonesia"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, -2, 0)),
			ViewCount:         892,
			LikeCount:         56,
			ShareCount:        23,
			ReadingTime:       7,
			IsFeatured:        false,
			IsSpotlight:       true,
			Priority:          8,
			Language:          "en",
			Version:           1,
			EstimatedTime:     7,
			AvgReadTime:       6.2,
			BounceRate:        0.41,
			CompletionRate:    0.65,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000003"),
			Title:          "Night Driving: Essential Tips for Safety",
			Slug:           "night-driving-essential-tips-for-safety",
			LeadParagraph:  "Night driving presents unique challenges. Learn how to stay safe when driving after dark.",
			Content: `<p>Night driving requires extra attention and different techniques than daytime driving.</p>

<h2>Check Your Lights</h2>
<p>Before driving at night, ensure all your lights are working properly. This includes headlights, tail lights, brake lights, and turn signals.</p>

<h2>Adjust Your Speed</h2>
<p>Reduce your speed to account for reduced visibility. The faster you drive, the less time you have to react.</p>

<h2>Use High Beams Wisely</h2>
<p>Use high beams on roads with minimal traffic. Switch to low beams when approaching other vehicles.</p>

<h2>Watch for Pedestrians</h2>
<p>Pedestrians can be harder to see at night. Be extra vigilant near crosswalks and populated areas.</p>`,
			Footer:           "Drive safely and arrive alive.",
			MetaTitle:         "Night Driving Safety Tips - Complete Guide",
			MetaDescription:   "Master night driving with these essential safety tips. Learn techniques for driving safely after dark.",
			MetaKeywords:      "night driving, driving safety, evening driving, nighttime driving tips",
			Tags:              pq.StringArray{"safety", "night-driving", "advanced", "driving-tips"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, -3, 0)),
			ViewCount:         654,
			LikeCount:         42,
			ShareCount:        18,
			ReadingTime:       4,
			IsFeatured:        false,
			IsSpotlight:       false,
			Priority:          5,
			Language:          "en",
			Version:           1,
			EstimatedTime:     4,
			AvgReadTime:       3.8,
			BounceRate:        0.38,
			CompletionRate:    0.72,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000004"),
			Title:          "How to Handle Wet Weather Driving",
			Slug:           "how-to-handle-wet-weather-driving",
			LeadParagraph:  "Rainy conditions require special driving techniques. Here's how to stay safe on wet roads.",
			Content: `<p>Indonesia's rainy season brings challenging driving conditions. Learn how to navigate safely.</p>

<h2>Reduce Speed</h2>
<p>Slow down when roads are wet. Hydroplaning becomes a risk at higher speeds.</p>

<h2>Increase Following Distance</h2>
<p>Wet roads mean longer braking distances. Keep more space between you and the vehicle ahead.</p>

<h2>Use Defoggers</h2>
<p>Keep your windshield clear using defoggers. Humidity can quickly fog up windows.</p>

<h2>Check Tire Condition</h2>
<p>Ensure your tires have adequate tread depth for proper water dispersion.</p>`,
			Footer:           "Stay safe during rainy season by following these guidelines.",
			MetaTitle:         "Wet Weather Driving Guide - Rain Safety Tips",
			MetaDescription:   "Learn essential techniques for driving safely in wet weather conditions. Complete guide for Indonesian drivers.",
			MetaKeywords:      "wet weather driving, rain driving, hydroplaning, driving safety",
			Tags:              pq.StringArray{"safety", "weather", "rain", "advanced-tips"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, -4, 0)),
			ViewCount:         445,
			LikeCount:         28,
			ShareCount:        12,
			ReadingTime:       5,
			IsFeatured:        false,
			IsSpotlight:       false,
			Priority:          4,
			Language:          "en",
			Version:           1,
			EstimatedTime:     5,
			AvgReadTime:       4.2,
			BounceRate:        0.45,
			CompletionRate:    0.68,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000005"),
			Title:          "The Benefits of Enrolling in a Professional Driving School",
			Slug:           "benefits-professional-driving-school",
			LeadParagraph:  "Professional driving schools offer advantages that self-learning simply cannot match.",
			Content: `<p>While some people learn to drive from family or friends, professional driving schools provide structured, comprehensive training.</p>

<h2>Structured Curriculum</h2>
<p>Professional schools follow a proven curriculum that covers all essential skills in a logical progression.</p>

<h2>Certified Instructors</h2>
<p>Learn from trained professionals who know exactly what examiners look for.</p>

<h2>Proper Equipment</h2>
<p>Driving schools maintain vehicles equipped with dual controls for safety.</p>

<h2>Higher Pass Rates</h2>
<p>Students from professional schools typically have higher first-time pass rates.</p>`,
			Footer:           "Invest in professional training for better driving skills.",
			MetaTitle:         "Why Choose Professional Driving School - Benefits Guide",
			MetaDescription:   "Discover the benefits of enrolling in a professional driving school vs self-learning.",
			MetaKeywords:      "driving school, professional training, driving lessons, SIM A preparation",
			Tags:              pq.StringArray{"driving-school", "tips", "beginners", "preparation"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, 0, -15)),
			ViewCount:         321,
			LikeCount:         24,
			ShareCount:        15,
			ReadingTime:       4,
			IsFeatured:        false,
			IsSpotlight:       false,
			Priority:          6,
			Language:          "en",
			Version:           1,
			EstimatedTime:     4,
			CTAText:           "View Our Packages",
			CTALink:           "/packages",
			AvgReadTime:       3.5,
			BounceRate:        0.35,
			CompletionRate:    0.75,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000006"),
			Title:          "Understanding Car Maintenance Basics",
			Slug:           "understanding-car-maintenance-basics",
			LeadParagraph:  "Every driver should know basic car maintenance. Keep your vehicle running smoothly with these essential tips.",
			Content: `<p>Basic car maintenance knowledge helps prevent breakdowns and extends your vehicle's lifespan.</p>

<h2>Check Fluid Levels</h2>
<p>Regularly check engine oil, coolant, brake fluid, and windshield washer fluid levels.</p>

<h2>Monitor Tire Pressure</h2>
<p>Proper tire pressure improves fuel efficiency and safety. Check monthly.</p>

<h2>Replace Wipers</h2>
<p>Replace windshield wipers every 6-12 months for clear visibility during rain.</p>

<h2>Schedule Regular Services</h2>
<p>Follow your car's service schedule for oil changes, filter replacements, and inspections.</p>`,
			Footer:           "A well-maintained car is a safer car.",
			MetaTitle:         "Car Maintenance Basics - Every Driver's Guide",
			MetaDescription:   "Learn essential car maintenance tips every driver should know. Keep your vehicle in top condition.",
			MetaKeywords:      "car maintenance, vehicle care, basic car tips, driving knowledge",
			Tags:              pq.StringArray{"maintenance", "car-care", "tips", "knowledge"},
			Status:            models.ArticleStatusPublished,
			PublishedAt:       ptrTime(time.Now().AddDate(0, 0, -7)),
			ViewCount:         278,
			LikeCount:         19,
			ShareCount:        8,
			ReadingTime:       6,
			IsFeatured:        false,
			IsSpotlight:       false,
			Priority:          3,
			Language:          "en",
			Version:           1,
			EstimatedTime:     6,
			AvgReadTime:       5.1,
			BounceRate:        0.42,
			CompletionRate:    0.62,
		},
		{
			ID:             uuid.MustParse("a0000001-0000-0000-0000-000000000007"),
			Title:          "Parking Techniques for Beginners",
			Slug:           "parking-techniques-for-beginners",
			LeadParagraph:  "Parking is often the biggest challenge for new drivers. Master these techniques to park with confidence.",
			Content: `<p>Good parking skills are essential for everyday driving. Practice these techniques regularly.</p>

<h2>Parallel Parking</h2>
<p>Find a space 1.5x your car's length. Use mirrors and turn indicators. Reverse slowly while turning the wheel.</p>

<h2>Perpendicular Parking</h2>
<p>Position your car parallel to the space. Reverse until aligned, then turn into the space.</p>

<h2>Angle Parking</h2>
<p>Approach at an angle. Use reference points to judge distance and steering angle.</p>

<h2>Tips for Success</h2>
<p>Always check mirrors and blind spots. Take your time - rushing leads to mistakes.</p>`,
			Footer:           "Practice makes perfect when it comes to parking.",
			MetaTitle:         "Parking Techniques for Beginners - Complete Guide",
			MetaDescription:   "Learn essential parking techniques for new drivers. Master parallel, perpendicular, and angle parking.",
			MetaKeywords:      "parking techniques, parallel parking, driving skills, beginners",
			Tags:              pq.StringArray{"parking", "skills", "beginners", "techniques"},
			Status:            models.ArticleStatusDraft,
			ViewCount:         0,
			LikeCount:         0,
			ShareCount:        0,
			ReadingTime:       5,
			IsFeatured:        false,
			IsSpotlight:       false,
			Priority:          2,
			Language:          "en",
			Version:           1,
			EstimatedTime:     5,
			AvgReadTime:       0,
			BounceRate:        0,
			CompletionRate:    0,
		},
	}

	// Create categories
	categories := []models.Category{
		{
			ID:          uuid.MustParse("c0000001-0000-0000-0000-000000000001"),
			Name:        "Driving Tips",
			Slug:        "driving-tips",
			Description: "Practical advice for better driving",
			Order:       1,
		},
		{
			ID:          uuid.MustParse("c0000001-0000-0000-0000-000000000002"),
			Name:        "Road Safety",
			Slug:        "road-safety",
			Description: "Stay safe on the roads",
			Order:       2,
		},
		{
			ID:          uuid.MustParse("c0000001-0000-0000-0000-000000000003"),
			Name:        "Beginners Guide",
			Slug:        "beginners-guide",
			Description: "Everything new drivers need to know",
			Order:       3,
		},
	}

	// Create tags
	tags := []models.Tag{
		{
			ID:          uuid.MustParse("b0000001-0000-0000-0000-000000000001"),
			Name:        "Tips",
			Slug:        "tips",
			Description: "Helpful driving tips",
			UsageCount:  5,
		},
		{
			ID:          uuid.MustParse("b0000001-0000-0000-0000-000000000002"),
			Name:        "Safety",
			Slug:        "safety",
			Description: "Road safety content",
			UsageCount:  4,
		},
		{
			ID:          uuid.MustParse("b0000001-0000-0000-0000-000000000003"),
			Name:        "Beginners",
			Slug:        "beginners",
			Description: "Content for new drivers",
			UsageCount:  3,
		},
	}

	// Create categories first
	for _, category := range categories {
		result := db.Where("id = ?", category.ID).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}
	}

	// Create tags
	for _, tag := range tags {
		result := db.Where("id = ?", tag.ID).FirstOrCreate(&tag)
		if result.Error != nil {
			return result.Error
		}
	}

	// Update articles with category IDs
	articles[0].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000001"))
	articles[1].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000002"))
	articles[2].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000002"))
	articles[3].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000002"))
	articles[4].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000003"))
	articles[5].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000001"))
	articles[6].CategoryID = ptrUUID(uuid.MustParse("c0000001-0000-0000-0000-000000000001"))

	// Set author ID (use a generic UUID for seeding)
	authorID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	for i := range articles {
		articles[i].AuthorID = authorID
		articles[i].CreatedAt = time.Now()
		articles[i].UpdatedAt = time.Now()
	}

	// Create articles
	for _, article := range articles {
		result := db.Where("id = ?", article.ID).FirstOrCreate(&article)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Helper functions
func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrUUID(u uuid.UUID) *uuid.UUID {
	return &u
}